package backup

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Metrics is a map of default MetricDescriptions for this namespace
var Metrics = map[string]*b.MetricDescription{
	"NumberOfBackupJobsCreated": {
		Help:          aws.String("The number of backup jobs that AWS Backup created."),
		OutputName:    aws.String("backup_number_of_backup_jobs_created"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfBackupJobsPending": {
		Help:          aws.String("The number of backup jobs about to run in AWS Backup."),
		OutputName:    aws.String("backup_number_of_backup_jobs_pending"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfBackupJobsRunning": {
		Help:          aws.String("The number of backup jobs currently running in AWS Backup."),
		OutputName:    aws.String("backup_number_of_backup_jobs_running"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfBackupJobsAborted": {
		Help:          aws.String("The number of user cancelled backup jobs."),
		OutputName:    aws.String("backup_number_of_backup_jobs_aborted"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfBackupJobsCompleted": {
		Help:          aws.String("NumberOfBackupJobsCompleted"),
		OutputName:    aws.String("backup_number_of_backup_jobs_completed"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfBackupJobsFailed": {
		Help:          aws.String("The number of backup jobs that AWS Backup scheduled but did not start. Often caused by scheduling a backup job during or 1 hour before a database resource or 4 hours before or during a Amazon FSx maintenance window or automated backup window and not using AWS Backup to perform continuous backup for point-in-time restores. See Point-in-Time Recovery for a list of supported services and instructions on how to use AWS Backup to take continuous backups, or reschedule your backup jobs."),
		OutputName:    aws.String("backup_number_of_backup_jobs_failed"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfBackupJobsExpired": {
		Help:          aws.String("The number of backup jobs that AWS Backup attempted to delete based on your backup retention lifecycle, but could not delete. You are billed for the storage that expired backups consume and should delete them manually."),
		OutputName:    aws.String("backup_number_of_backup_jobs_expired"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfCopyJobsCreated": {
		Help:          aws.String("The number of cross-account and cross-Region copy jobs that AWS Backup created."),
		OutputName:    aws.String("backup_number_of_copy_jobs_created"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfCopyJobsRunning": {
		Help:          aws.String("The number of cross-account and cross-Region copy jobs currently running in AWS Backup."),
		OutputName:    aws.String("backup_number_of_copy_jobs_running"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfCopyJobsCompleted": {
		Help:          aws.String("The number of cross-account and cross-Region copy jobs that AWS Backup finished."),
		OutputName:    aws.String("backup_number_of_copy_jobs_completed"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfCopyJobsFailed": {
		Help:          aws.String("The number of cross-account and cross-Region copy jobs that AWS Backup attempted but could not complete."),
		OutputName:    aws.String("backup_number_of_copy_jobs_failed"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfRestoreJobsPending": {
		Help:          aws.String("The number of restore jobs about to run in AWS Backup."),
		OutputName:    aws.String("backup_number_of_restore_jobs_pending"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfRestoreJobsRunning": {
		Help:          aws.String("The number of restore jobs currently running in AWS Backup."),
		OutputName:    aws.String("backup_number_of_restore_jobs_running"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfRestoreJobsCompleted": {
		Help:          aws.String("The number of restore jobs that AWS Backup finished."),
		OutputName:    aws.String("backup_number_of_restore_jobs_completed"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfRestoreJobsFailed": {
		Help:       aws.String("The number of restore jobs that AWS Backup attempted but could not complete."),
		OutputName: aws.String("backup_number_of_restore_jobs_failed"),
		Statistic:  h.StringPointers("Average"),
		Kind:       aws.String(b.CLOUDWATCH_KIND),

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfRecoveryPointsCompleted": {
		Help:          aws.String("The number of recovery points that AWS Backup created."),
		OutputName:    aws.String("backup_number_of_recovery_points_completed"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfRecoveryPointsPartial": {
		Help:          aws.String("The number of recovery points that AWS Backup started to create but could not finish. AWS retries the process later, but because the retry occurs at the later time, it retains the partial recovery point."),
		OutputName:    aws.String("backup_number_of_recovery_points_partial"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfRecoveryPointsExpired": {
		Help:          aws.String("The number of recovery points that AWS Backup attempted to delete based on your backup retention lifecycle, but could not delete. You are billed for the storage that expired backups consume and should delete them manually."),
		OutputName:    aws.String("backup_mumber_of_recovery_points_expired"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfRecoveryPointsDeleting": {
		Help:          aws.String("The number of recovery points that AWS Backup is deleting."),
		OutputName:    aws.String("backup_number_of_recovery_points_deleting"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
	"NumberOfRecoveryPointsCold": {
		Help:          aws.String("The number of recovery points that AWS Backup tiered to cold storage."),
		OutputName:    aws.String("backup_number_of_recovery_points_cold"),
		Statistic:     h.StringPointers("Average"),
		Kind:          aws.String(b.CLOUDWATCH_KIND),
		PeriodSeconds: 60 * 5,
		RangeSeconds:  60 * 60 * 24,

		Dimensions: []*cloudwatch.Dimension{},
	},
}
