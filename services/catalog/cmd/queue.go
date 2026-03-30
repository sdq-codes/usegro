package cmd

//
//import (
//	"context"
//	"fmt"
//	"os"
//	"os/signal"
//	"sync"
//	"syscall"
//
//	"github.com/google/uuid"
//	"github.com/usegro/services/catalog/internal/helper/queue"
//	"github.com/usegro/services/catalog/internal/logger"
//	"github.com/spf13/cobra"
//	"go.uber.org/zap"
//)
//
//func init() {
//	rootCmd.AddGroup(&cobra.Group{ID: "queue", Title: "Queue:"})
//	rootCmd.AddCommand(
//		queueWorkCommand,
//		queueClearCommand,
//		queueFlushCommand,
//		queueForgetCommand,
//		queueRetryCommand,
//		queueRestoreCommand,
//	)
//
//	queueWorkCommand.Flags().StringP("queue", "q", "default", "(optional) queue name. for example: -q emails")
//	queueWorkCommand.Flags().IntP("worker", "w", 1, "(optional) The number of worker goroutines to run. for example: -w 2")
//	queueWorkCommand.Example = "  queue:work"
//	queueWorkCommand.Example += "\n  queue:work -w 2"
//	queueWorkCommand.Example += "\n  queue:work -q emails -w 2"
//
//	queueRetryCommand.Flags().StringP("queue", "q", "default", "(optional) queue name. for example: -q emails")
//	queueRetryCommand.Flags().StringP("id", "i", "", "(optional) job id. for example: --id df6df3af-d53d-49c2-bd50-80ba1d32b17b")
//	queueRetryCommand.Example = `  queue:retry -q emails -i df6df3af-d53d-49c2-bd50-80ba1d32b17b`
//
//	queueClearCommand.Flags().StringP("queue", "q", "default", "(optional) queue name. for example: -q emails")
//	queueClearCommand.Flags().BoolP("all", "a", false, "(optional) force delete all jobs for all queues.")
//	queueClearCommand.Example = "  queue:clear"
//	queueClearCommand.Example += "\n  queue:clear -q emails"
//	queueClearCommand.Example += "\n  queue:clear -a"
//
//	queueFlushCommand.Flags().StringP("queue", "q", "default", "(optional) queue name. for example: -q emails")
//	queueFlushCommand.Flags().BoolP("all", "a", false, "(optional) force delete all failed_jobs for all queues.")
//	queueFlushCommand.Example = "  queue:flush"
//	queueFlushCommand.Example += "\n  queue:flush -q emails"
//	queueFlushCommand.Example += "\n  queue:flush -a"
//
//	queueForgetCommand.Flags().StringP("id", "i", "", "(optional) job id. for example: --id df6df3af-d53d-49c2-bd50-80ba1d32b17b")
//	queueForgetCommand.MarkFlagRequired("id")
//	queueForgetCommand.Example = `  queue:forget -i df6df3af-d53d-49c2-bd50-80ba1d32b17b`
//
//	queueRestoreCommand.Flags().StringP("queue", "q", "default", "(optional) queue name. for example: -q emails")
//	queueRestoreCommand.Example = "  queue:restore"
//	queueRestoreCommand.Example += "\n  queue:restore -q emails"
//}
