package builder

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sarcasticadmin/sshcb/logs"
)

type SSHConfigOptions struct {
	Username     string
	Filepath     string
	BastionHost  string
	PrivateOnly  bool
	IdentityFile string
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
			return IncrementID(name+"-2", instanceDict)
		}
	}
	return name
}
func BuildInstanceList(reservations []*ec2.Reservation) map[string]InstanceInfo {
	instances := make(map[string]InstanceInfo)
	for idx, res := range reservations {
		logs.INFO.Println("  > Reservation Id", *res.ReservationId, " Num Instances: ", len(res.Instances))
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
			var ip string
			if inst.PublicIpAddress == "" && sshConfig.PrivateOnly == false {
				logs.WARN.Printf("Cannot find public IP for %s, skipping since bastion not set...\n", inst.InstanceID)
				continue
			} else if sshConfig.PrivateOnly == true {
				ip = inst.PrivateIpAddress
			} else {
				ip = inst.PublicIpAddress
			}
			s += fmt.Sprintf("# %s\nHost %s\n\tHostname  %s\n\tUser  %s\n",
				inst.InstanceID,
				name,
				ip,
				sshConfig.Username)
		}

		if sshConfig.IdentityFile != "" {
			s += fmt.Sprintf("\tIdentityFile %s\n",
				sshConfig.IdentityFile)
		}

	}
	newconfig, err := f.WriteString(s)
	logs.INFO.Printf("wrote %d bytes\n", newconfig)
	logs.FEEDBACK.Printf("Created ssh config at %s\n", sshConfig.Filepath)
	f.Sync()
}

func GetSession(profile string, region string) *ec2.EC2 {
	var ec2svc *ec2.EC2
	if profile == "" {
		ec2svc = ec2.New(session.New(&aws.Config{Region: aws.String(region)}))
	} else {
		ec2svc = ec2.New(session.New(&aws.Config{Region: aws.String(region),
			Credentials: credentials.NewSharedCredentials("", profile)}))
	}
	return ec2svc
}

func GetReservs(tags map[string]string, ec2svc *ec2.EC2) *ec2.DescribeInstancesOutput {
	logs.INFO.Println(tags)
	//fmt.Println(tags)
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

	//logs.INFO.Printf(filters)
	params := &ec2.DescribeInstancesInput{
		Filters: filters,
	}
	resp, err := ec2svc.DescribeInstances(params)
	if err != nil {
		logs.FATAL.Printf("there was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}
	return resp
}
