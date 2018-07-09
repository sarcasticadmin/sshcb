package builder

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"os"
	"strconv"
	"strings"
)

type SSHConfigOptions struct {
	Username    string
	Filepath    string
	BastionHost string
}

type InstanceInfo struct {
	InstanceID       string
	InstanceName     string
	PublicIpAddress  string
	PrivateIpAddress string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func IncrementID(name string, instanceDict map[string]InstanceInfo) string {
	if _, ok := instanceDict[name]; ok {
		testlist := strings.Split(name, "-")
		endval := testlist[len(testlist)-1]
		ival, err := strconv.Atoi(endval)
		if err == nil {
			ival = ival + 1
			testlist[len(testlist)-1] = strconv.Itoa(ival)
			return IncrementID(strings.Join(testlist, "-"), instanceDict)
		} else {
			return name + "-2"
		}
	}
	return name
}
func BuildInstanceList(reservations []*ec2.Reservation) map[string]InstanceInfo {
	instances := make(map[string]InstanceInfo)
	for idx, res := range reservations {
		fmt.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
		for _, inst := range reservations[idx].Instances {
			aID := InstanceInfo{
				InstanceID:       *inst.InstanceId,
				PrivateIpAddress: *inst.PrivateIpAddress,
			}
			if inst.PublicIpAddress != nil {
				aID.PublicIpAddress = *inst.PublicIpAddress
			}
			for _, tag := range inst.Tags {
				if *tag.Key == "Name" {
					aID.InstanceName = strings.ToLower(strings.Replace(*tag.Value, " ", "-", -1))
					break
				}
			}

			instances[IncrementID(aID.InstanceName, instances)] = (aID)

		}
	}
	return instances

}
func WriteSSHConfig(instanceList map[string]InstanceInfo, sshConfig SSHConfigOptions) {
	f, err := os.Create(sshConfig.Filepath)
	check(err)
	defer f.Close()
	s := ""
	if sshConfig.BastionHost != "" {
		s += fmt.Sprintf("Host %s\n\tHostname  %s\n\tUser  %s\n", sshConfig.BastionHost, sshConfig.BastionHost, sshConfig.Username)
	}

	for name, inst := range instanceList {
		if sshConfig.BastionHost != "" {
			s += fmt.Sprintf("# %s\nHost %s\n\tHostname  %s\n\tUser  %s\n\tProxyCommand ssh -F %s -W %%h:%%p %s\n",
				inst.InstanceID,
				name,
				inst.PrivateIpAddress,
				sshConfig.Username,
				sshConfig.Filepath,
				sshConfig.BastionHost)
		} else {
			if inst.PublicIpAddress == "" {
				fmt.Printf("Cannot find public IP for %s, skipping since bastion not set...\n", inst.InstanceID)
				continue
			}
			s += fmt.Sprintf("# %s\nHost %s\n\tHostname  %s\n\tUser  %s\n",
				inst.InstanceID,
				name,
				inst.PublicIpAddress,
				sshConfig.Username)
		}

	}
	newconfig, err := f.WriteString(s)
	fmt.Printf("wrote %d bytes\n", newconfig)
	f.Sync()
}

func GetSession(profile string, region string) *ec2.EC2 {
	ec2svc := ec2.New(session.New(&aws.Config{Region: aws.String(region),
		Credentials: credentials.NewSharedCredentials("", profile)}))
	return ec2svc
}

func GetReservs(tags map[string]string, ec2svc *ec2.EC2) *ec2.DescribeInstancesOutput {
	fmt.Println(tags)
	filters := []*ec2.Filter{}

	// We only care about running/pending instances
	filters = append(filters, &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []*string{aws.String("running"), aws.String("pending")},
	})

	// If any tags are present for the filtering include them as well
	for k, v := range tags {
		filters = append(filters, &ec2.Filter{
			Name:   aws.String(fmt.Sprintf("tag:%s", k)),
			Values: []*string{aws.String(v)},
		})
	}

	fmt.Println(filters)
	params := &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	/*
		params := &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("instance-state-name"),
					Values: []*string{aws.String("running"), aws.String("pending")},
				},
			},
		}
	*/
	resp, err := ec2svc.DescribeInstances(params)
	if err != nil {
		fmt.Println("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}
	return resp
}