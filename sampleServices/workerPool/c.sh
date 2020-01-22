rm -R .fossa.yml .git artifacts assets builds dev docs *.go Gopkg.lock Gopkg.toml logs Makefile manifests README.md sonarcloud.sh sonar-project.properties volumes vendor* worker 2> /dev/null
mkdir logs builds artifacts
chown jscharber:jscharber logs builds artifacts
