package base

var namespaces = []string{
	"AWS/ELB",
	"AWS/NetworkELB",
	"AWS/ApplicationELB",
	"AWS/RDS",
	"AWS/ElastiCache",
	"AWS/NATGateway",
	"AWS/EC2",
	"AWS/S3",
	"AWS/SQS",
	"AWS/VPC",
}

// GetNamespaces returns a list of AWS namespaces which are configured for this exporter
func GetNamespaces() []string {
	return namespaces
}
