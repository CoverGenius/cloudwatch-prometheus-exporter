package backup

import (
	b "github.com/CoverGenius/cloudwatch-prometheus-exporter/base"
	h "github.com/CoverGenius/cloudwatch-prometheus-exporter/helpers"
	log "github.com/sirupsen/logrus"

	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/backup"
)

func createResourceDescription(nd *b.NamespaceDescription, tags *string, vn *string) (*b.ResourceDescription, error) {
	rd := b.ResourceDescription{}
	dd := []*b.DimensionDescription{
		{
			Name:  aws.String("BackupVaultName"),
			Value: vn,
		},
	}
	err := rd.BuildDimensions(dd)
	h.LogIfError(err)
	rd.ID = vn
	rd.Name = vn
	rd.Type = aws.String("backup")
	rd.Parent = nd
	rd.Tags = tags

	return &rd, nil
}

// CreateResourceList fetches a list of all Backup resources in the region
func CreateResourceList(nd *b.NamespaceDescription, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Debug("Creating Backup resource list ...")
	session := backup.New(nd.Parent.Session)
	input := backup.ListBackupVaultsInput{}
	result, err := session.ListBackupVaults(&input)
	h.LogIfError(err)

	var w sync.WaitGroup
	w.Add(len(result.BackupVaultList))

	ch := make(chan *b.ResourceDescription, len(result.BackupVaultList))
	for _, vault := range result.BackupVaultList {
		go func(vlm *backup.VaultListMember, wg *sync.WaitGroup) {
			defer wg.Done()
			input := backup.ListTagsInput{
				ResourceArn: vlm.BackupVaultArn,
			}
			tags, err := session.ListTags(&input)
			h.LogIfError(err)

			tl, found := nd.Parent.TagsFound(tags)
			ts := b.TagsToString(tl)

			if found {
				if r, err := createResourceDescription(nd, ts, vlm.BackupVaultName); err == nil {
					ch <- r
				}
				h.LogIfError(err)
			}
		}(vault, &w)
	}
	w.Wait()
	close(ch)

	resources := []*b.ResourceDescription{}
	for r := range ch {
		resources = append(resources, r)
	}
	nd.Mutex.Lock()
	nd.Resources = resources
	nd.Mutex.Unlock()

}
