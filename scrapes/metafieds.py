"""
Shopify Taxonomy Metafields Scraper
=====================================
Fetches all category metafields (product attributes) from Shopify's
public product taxonomy. Uses the Shopify Admin GraphQL API which returns
the full taxonomy tree with attributes per category — no guesswork about
file structure needed.

TWO MODES:
  1. GraphQL API (recommended) — requires a store domain + Admin API token
  2. GitHub JSON files (fallback) — no token needed, parses categories.json
     which at 61.9MB contains the full tree including per-category attributes

Output: shopify_metafields.json

Usage:
    pip install requests

    # Mode 1: GraphQL (most reliable)
    python scrape_shopify_metafields.py --store yourstore.myshopify.com --token shpat_xxx

    # Mode 2: GitHub public files (no token)
    python scrape_shopify_metafields.py --github
"""

import argparse
import json
import sys
import time
import requests

# ─── Constants ────────────────────────────────────────────────────────────────

GITHUB_BASE = "https://raw.githubusercontent.com/Shopify/product-taxonomy/main/dist/en"
OUTPUT_FILE = "shopify_metafields.json"

# GraphQL query — fetches taxonomy categories with their attribute definitions
TAXONOMY_QUERY = """
query GetTaxonomyCategories($after: String) {
  taxonomy {
    categories(first: 250, after: $after) {
      pageInfo {
        hasNextPage
        endCursor
      }
      edges {
        node {
          id
          name
          fullName
          level
          isLeaf
          isRoot
          parentId
          attributes {
            id
            name
            handle
            values(first: 250) {
              id
              name
              handle
            }
          }
        }
      }
    }
  }
}
"""


# ─── Mode 1: Shopify Admin GraphQL ────────────────────────────────────────────

def fetch_via_graphql(store: str, token: str) -> dict:
    """Paginate through all taxonomy categories via Shopify Admin GraphQL API."""
    endpoint = f"https://{store}/admin/api/2025-01/graphql.json"
    headers = {
        "Content-Type": "application/json",
        "X-Shopify-Access-Token": token,
    }

    all_categories = []
    cursor = None
    page = 1

    print(f"\n  Endpoint: {endpoint}")

    while True:
        print(f"  Fetching page {page} (cursor: {cursor or 'start'})...")
        payload = {"query": TAXONOMY_QUERY, "variables": {"after": cursor}}
        resp = requests.post(endpoint, json=payload, headers=headers, timeout=30)
        resp.raise_for_status()

        data = resp.json()
        if "errors" in data:
            print(f"\n❌ GraphQL errors: {json.dumps(data['errors'], indent=2)}")
            sys.exit(1)

        cats = data["data"]["taxonomy"]["categories"]
        edges = cats["edges"]
        all_categories.extend(e["node"] for e in edges)

        page_info = cats["pageInfo"]
        if not page_info["hasNextPage"]:
            break

        cursor = page_info["endCursor"]
        page += 1
        time.sleep(0.3)  # be polite to the API

    return build_output_graphql(all_categories)


# ─── Mode 2: GitHub public JSON files ─────────────────────────────────────────

def fetch_github_json(path: str) -> any:
    url = f"{GITHUB_BASE}/{path}"
    print(f"  GET {url}")
    resp = requests.get(url, timeout=120)
    resp.raise_for_status()
    return resp.json()


def fetch_via_github() -> dict:
    """
    Parse Shopify's public taxonomy dist files from GitHub.

    categories.json (61.9 MB) contains the full category tree.
    Each category node looks like:
      {
        "id": "gid://shopify/TaxonomyCategory/aa",
        "name": "Apparel & Accessories",
        "full_name": "Apparel & Accessories",
        "level": 0,
        "children_ids": [...],
        "attributes": [
          { "id": "gid://shopify/TaxonomyAttribute/1", "handle": "color", "name": "Color" },
          ...
        ]
      }

    attributes.json contains the full attribute definitions including all values.
    """
    print("\n[1/2] Downloading attributes.json...")
    attrs_raw = fetch_github_json("attributes.json")

    print("[2/2] Downloading categories.json (~62 MB) — this may take a minute...")
    cats_raw = fetch_github_json("categories.json")

    # ── Extract attributes list ──────────────────────────────────────────────
    def extract_list(data):
        if isinstance(data, list):
            return data
        if isinstance(data, dict):
            for key in ("attributes", "categories", "data"):
                if isinstance(data.get(key), list):
                    return data[key]
            # fallback: first list value found
            for v in data.values():
                if isinstance(v, list):
                    return v
        return []

    attrs_list = [a for a in extract_list(attrs_raw) if isinstance(a, dict)]
    print(f"\n  ✓ {len(attrs_list)} attributes loaded")

    # Build lookup by both id and handle
    attr_by_id     = {a["id"]: a     for a in attrs_list if a.get("id")}
    attr_by_handle = {a["handle"]: a for a in attrs_list if a.get("handle")}

    # ── Extract categories list ──────────────────────────────────────────────
    # categories.json structure: { "verticals": [ { "categories": [...] } ] }
    # Each vertical's "categories" is a flat list of all categories in that vertical
    cats_list = []
    for vertical in cats_raw.get("verticals", []):
        cats_list.extend(c for c in vertical.get("categories", []) if isinstance(c, dict))

    if not cats_list:
        # fallback for unexpected structure
        cats_list = [c for c in extract_list(cats_raw) if isinstance(c, dict)]

    print(f"  ✓ {len(cats_list)} categories loaded")

    # ── Resolve attributes per category ─────────────────────────────────────
    normalized_cats = []
    for cat in cats_list:
        raw_attrs = (
            cat.get("attributes")
            or cat.get("attribute_ids")
            or cat.get("attribute_handles")
            or []
        )

        resolved = []
        for entry in raw_attrs:
            if isinstance(entry, dict):
                # May already have id/handle — enrich from full lookup
                aid = entry.get("id")
                ahandle = entry.get("handle")
                full = attr_by_id.get(aid) or attr_by_handle.get(ahandle) or entry
                resolved.append(_format_attr(full))

            elif isinstance(entry, str):
                # String is either a GID or a handle
                full = attr_by_id.get(entry) or attr_by_handle.get(entry)
                if full:
                    resolved.append(_format_attr(full))
                else:
                    # Unknown — preserve it so nothing is silently dropped
                    resolved.append({"id": entry, "name": entry, "handle": entry, "value_count": 0, "values": []})

        normalized_cats.append({
            "id":        cat.get("id"),
            "name":      cat.get("name"),
            "full_name": cat.get("full_name") or cat.get("fullName"),
            "level":     cat.get("level"),
            "is_leaf":   cat.get("is_leaf") or cat.get("isLeaf"),
            "attributes": resolved,
        })

    return build_output_github(attrs_list, normalized_cats)


def _format_attr(attr: dict) -> dict:
    values = attr.get("values") or []
    # values may be a list of dicts or a GraphQL connection {"nodes": [...]}
    if isinstance(values, dict):
        values = values.get("nodes", [])
    return {
        "id":          attr.get("id"),
        "name":        attr.get("name"),
        "handle":      attr.get("handle"),
        "value_count": len(values),
        "values": [
            {"id": v.get("id"), "name": v.get("name"), "handle": v.get("handle")}
            for v in values if isinstance(v, dict)
        ],
    }


# ─── Output builders ──────────────────────────────────────────────────────────

def build_output_graphql(graphql_categories: list) -> dict:
    # Collect unique attributes across all categories
    seen = {}
    for cat in graphql_categories:
        for attr in cat.get("attributes", []):
            aid = attr.get("id")
            if aid and aid not in seen:
                seen[aid] = attr

    return {
        "source": "graphql",
        "total_categories":         len(graphql_categories),
        "total_unique_attributes":  len(seen),
        "attributes": [_format_attr(a) for a in seen.values()],
        "categories": [
            {
                "id":        c.get("id"),
                "name":      c.get("name"),
                "full_name": c.get("fullName") or c.get("full_name"),
                "level":     c.get("level"),
                "is_leaf":   c.get("isLeaf") or c.get("is_leaf"),
                "attributes": [_format_attr(a) for a in c.get("attributes", [])],
            }
            for c in graphql_categories
        ],
    }


def build_output_github(attrs_list: list, categories: list) -> dict:
    return {
        "source": "github",
        "total_categories":        len(categories),
        "total_unique_attributes": len(attrs_list),
        "attributes": [_format_attr(a) for a in attrs_list],
        "categories": categories,
    }


# ─── Main ─────────────────────────────────────────────────────────────────────

def main():
    parser = argparse.ArgumentParser(description="Scrape Shopify taxonomy metafields")
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument("--github", action="store_true",
                       help="Fetch from Shopify's public GitHub repo (no token needed)")
    group.add_argument("--store", metavar="DOMAIN",
                       help="Shopify store domain e.g. yourstore.myshopify.com")
    parser.add_argument("--token", metavar="TOKEN",
                        help="Shopify Admin API access token (required with --store)")
    args = parser.parse_args()

    print("\n🛍  Shopify Taxonomy Metafields Scraper\n" + "=" * 42)

    if args.github:
        print("\nMode: GitHub public files")
        output = fetch_via_github()
    else:
        if not args.token:
            print("❌ --token is required when using --store")
            sys.exit(1)
        print(f"\nMode: Shopify Admin GraphQL ({args.store})")
        output = fetch_via_graphql(args.store, args.token)

    with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
        json.dump(output, f, indent=2, ensure_ascii=False)

    cats_with_attrs = sum(1 for c in output["categories"] if c.get("attributes"))

    print("\n✅ Done!")
    print(f"   Total categories:           {output['total_categories']}")
    print(f"   Categories with metafields: {cats_with_attrs}")
    print(f"   Unique attributes:          {output['total_unique_attributes']}")
    print(f"\n📄 Saved to: {OUTPUT_FILE}\n")


if __name__ == "__main__":
    main()