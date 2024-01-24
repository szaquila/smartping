#/bin/bash

# baseDirForScriptSelf=$(cd "$(dirname "$0")"; pwd)
# echo "full path to currently executed script is : ${baseDirForScriptSelf}"
PC=`ps -ef|grep "bin/smartping"|grep -v grep|wc -l`
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ];
do
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
cd "$DIR"
# cd "$(dirname "$0")"
if [ $PC == 0 ]; then
  ./control start
# else
# ./control restart
fi

