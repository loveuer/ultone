#/bin/bash

VERSION="v$(date +'%y.%m.%d')-r1"
echo "version: $VERSION"

docker build -t repository.umisen.com/{project_folder}/{project_name}:$VERSION -f Dockerfile .
docker push repository.umisen.com/{project_folder}/{project_name}:$VERSION
docker start -d --name {your_project_name} --restart unless-stopped -p xx_port:xx_port -v xx_path:xx_path repository.umisen.com/{project_folder}/{project_name}:$VERSION