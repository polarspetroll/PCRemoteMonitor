[InstallServer]
cd Server
go build -o PCRemoteMonitorService
sudo mkdir /usr/bin/PCRemoteMonitor
sudo mv PCRemoteMonitorService  /usr/bin/PCRemoteMonitor
sudo cp config.json /usr/bin/PCRemoteMonitor
sudo cp PCRemoteMonitor.service /etc/systemd/system/
sudo systemctl start PCRemoteMonitor
sudo systemctl enable PCRemoteMonitor

[buildAndroidApp]
cd App
fyne package -os android
