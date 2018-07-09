# SSH Config Builder

Used to generate an `ssh_config` from cloud resources (only aws currently)
> Note: Currently experimental

## Build
Assuming GOPATH is correctly configured:
```
go get github.com/sarcasticadmin/sshcb
cd $GOPATH/src/github.com/sarcasticadmin/sshcb
make build
```

## Examples

Create ssh_config from all `ec2` instances in `us-west-2` with user `ubuntu`:
```
sshcb -r us-west-2 -u ubuntu -c ~/.ssh/config.us-west-2
ssh -F ~/.ssh/config.us-west-2 coolinstance
```

Create ssh_config with aws profile `env1` and only with ec2 instances tagged `env:prod`:
```
sshcb -r us-west-2 --tags 'env:prod'
```
