set -e

Ver=2.0-dev
BuildT=$(date -u +'%Y-%m-%dT%H:%M:%SZ')
GitCommit=$(git rev-parse --short HEAD)
echo "\$Ver:    "     $Ver        ;
echo "\$BuildT: "     $BuildT     ;
echo "\$GitCommit: "  $GitCommit  ;

Dir=tmp/release
echo "Dir: "$Dir
rm -rf $Dir

go build -o /app/main main.go && ls -alh /app
