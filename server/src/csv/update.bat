SET CUR_PATH=%~dp0
SET CONFIG_NEW="E:/md_config/Magic/config_new"
SET CONFIG_CSV="E:/md_config/Magic/csv_new"
SET SRC_CSV=%~dp0\csv
SET SRC=%~dp0

REM gox -osarch "linux/amd64" -output "src"
TortoiseProc.exe /command:update /path:%CONFIG_NEW% /closeonend:3
echo update config ok!
TortoiseProc.exe /command:update /path:%CONFIG_CSV% /closeonend:3
echo update csv ok!
cd %CONFIG_CSV%
excel2csv.exe
svn add . --no-ignore --force
svn ci -m "update csv"
TortoiseProc.exe /command:update /path:%SRC_CSV% /closeonend:3
echo csv.zip ok!
cd %SRC%
REM robocopy %CONFIG_CSV% *.csv %SRC_CSV%
robocopy %CONFIG_CSV% %SRC_CSV% /S /ZB /R:3 /W:10 /V /MT:16
echo excel2csv ok
cd %SRC_CSV%
svn add . --no-ignore --force
svn ci -m "update config"
echo update csv done!
TortoiseProc.exe /command:update /path:%SRC% /closeonend:3
echo update src
echo start build exe!
cd %SRC%
del %SRC%src
SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o src
echo build ok!
cd %SRC%
del %SRC%csv.zip
%SRC%xupload.exe
echo upload done!

