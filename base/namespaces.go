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
}

func GetNamespaces() []string {
	return namespaces
}
