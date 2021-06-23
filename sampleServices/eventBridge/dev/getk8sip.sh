
#!/bin/bash

microk8s.config | grep server | sed -r "s/.*\/\/(.*):.*$/\1/"
