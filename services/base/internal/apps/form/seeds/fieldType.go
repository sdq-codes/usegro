package seeds

var FieldTypesSeed = []map[string]interface{}{
	{
		"name":        "Short Text",
		"description": "A single line of text input",
		"configs": []map[string]string{
			{"key": "placeholder", "valueType": "string", "description": "Placeholder text"},
			{"key": "minLength", "valueType": "int", "description": "Minimum characters"},
			{"key": "maxLength", "valueType": "int", "description": "Maximum characters"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "This field is required"},
			{"key": "regex", "rule": "string", "message": "Value must match pattern"},
		},
	},
	{
		"name":        "Long Text",
		"description": "Multi-line text area",
		"configs": []map[string]string{
			{"key": "placeholder", "valueType": "string", "description": "Placeholder text"},
			{"key": "minLength", "valueType": "int", "description": "Minimum characters"},
			{"key": "maxLength", "valueType": "int", "description": "Maximum characters"},
			{"key": "rows", "valueType": "int", "description": "Number of visible rows"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "This field is required"},
		},
	},
	{
		"name":        "Radio Button",
		"description": "User selects one option",
		"configs": []map[string]string{
			{"key": "options", "valueType": "array", "description": "List of available options"},
			{"key": "randomize", "valueType": "boolean", "description": "Randomize option order"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "An option must be selected"},
		},
	},
	{
		"name":        "Checkbox",
		"description": "User selects one or more options",
		"configs": []map[string]string{
			{"key": "options", "valueType": "array", "description": "List of available options"},
			{"key": "minSelections", "valueType": "int", "description": "Minimum selections allowed"},
			{"key": "maxSelections", "valueType": "int", "description": "Maximum selections allowed"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "At least one option is required"},
		},
	},
	{
		"name":        "Dropdown",
		"description": "User selects from a dropdown",
		"configs": []map[string]string{
			{"key": "options", "valueType": "array", "description": "List of dropdown options"},
			{"key": "placeholder", "valueType": "string", "description": "Placeholder text"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Selection is required"},
		},
	},
	{
		"name":        "File Upload",
		"description": "User uploads a file",
		"configs": []map[string]string{
			{"key": "allowedTypes", "valueType": "array", "description": "Allowed MIME types"},
			{"key": "maxSizeMB", "valueType": "int", "description": "Maximum file size in MB"},
			{"key": "multiple", "valueType": "boolean", "description": "Allow multiple files"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "File upload required"},
		},
	},
	{
		"name":        "Linear Scale",
		"description": "Numeric scale input",
		"configs": []map[string]string{
			{"key": "min", "valueType": "int", "description": "Minimum value"},
			{"key": "max", "valueType": "int", "description": "Maximum value"},
			{"key": "step", "valueType": "int", "description": "Step size"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Scale selection required"},
		},
	},
	{
		"name":        "Rating",
		"description": "Star rating input",
		"configs": []map[string]string{
			{"key": "maxStars", "valueType": "int", "description": "Maximum number of stars"},
			{"key": "icon", "valueType": "string", "description": "Custom icon for rating"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Rating required"},
		},
	},
	{
		"name":        "Date",
		"description": "Date picker",
		"configs": []map[string]string{
			{"key": "minDate", "valueType": "date", "description": "Earliest selectable date"},
			{"key": "maxDate", "valueType": "date", "description": "Latest selectable date"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Date required"},
		},
	},
	{
		"name":        "Time",
		"description": "Time picker",
		"configs": []map[string]string{
			{"key": "format", "valueType": "string", "description": "Display format (e.g. HH:mm)"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Time required"},
		},
	},
	{
		"name":        "Datetime",
		"description": "Date & time picker",
		"configs": []map[string]string{
			{"key": "minDateTime", "valueType": "datetime", "description": "Earliest date-time"},
			{"key": "maxDateTime", "valueType": "datetime", "description": "Latest date-time"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Datetime required"},
		},
	},
	{
		"name":        "Picture Choice",
		"description": "User selects from images",
		"configs": []map[string]string{
			{"key": "images", "valueType": "array", "description": "List of image URLs"},
			{"key": "multiple", "valueType": "boolean", "description": "Allow multiple selections"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Image choice required"},
		},
	},
	{
		"name":        "Yes/No",
		"description": "Boolean toggle",
		"configs": []map[string]string{
			{"key": "labels", "valueType": "array", "description": "Labels for yes/no"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Yes/No required"},
		},
	},
	{
		"name":        "Number",
		"description": "Numeric input",
		"configs": []map[string]string{
			{"key": "min", "valueType": "int", "description": "Minimum value"},
			{"key": "max", "valueType": "int", "description": "Maximum value"},
			{"key": "step", "valueType": "int", "description": "Increment step"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Number required"},
		},
	},
	{
		"name":        "Email",
		"description": "Valid email input",
		"configs": []map[string]string{
			{"key": "placeholder", "valueType": "string", "description": "Placeholder text"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Email is required"},
			{"key": "regex", "rule": `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, "message": "Must be a valid email"},
		},
	},
	{
		"name":        "Phone Number",
		"description": "Phone number input",
		"configs": []map[string]string{
			{"key": "placeholder", "valueType": "string", "description": "Placeholder text"},
			{"key": "countryCode", "valueType": "string", "description": "Default country code"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Phone number required"},
			{"key": "regex", "rule": `^\+?[1-9]\d{1,14}$`, "message": "Must be a valid phone number"},
		},
	},
	{
		"name":        "Address",
		"description": "Structured address input",
		"configs": []map[string]string{
			{"key": "fields", "valueType": "array", "description": "Street, City, State, Zip, Country"},
		},
		"validations": []map[string]string{
			{"key": "required", "rule": "boolean", "message": "Address required"},
		},
	},
	{
		"name":        "Link",
		"description": "URL input",
		"configs": []map[string]string{
			{"key": "placeholder", "valueType": "string", "description": "Placeholder text"},
		},
		"validations": []map[string]string{
			{"key": "regex", "rule": `^https?://[\w\-]+(\.[\w\-]+)+[/#?]?.*$`, "message": "Must be a valid URL"},
		},
	},
	{
		"name":        "Country",
		"description": "User selects country from a dropdown",
		"configs":     []map[string]string{},
		"validations": []map[string]string{},
	},
	{
		"name":        "State",
		"description": "User selects state from a dropdown",
		"configs":     []map[string]string{},
		"validations": []map[string]string{},
	},
	{
		"name":        "City",
		"description": "User selects city from a dropdown",
		"configs":     []map[string]string{},
		"validations": []map[string]string{},
	},
	{
		"name":        "Tags",
		"description": "Select tags",
		"configs":     []map[string]string{},
		"validations": []map[string]string{},
	},
}
