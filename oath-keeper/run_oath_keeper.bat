@echo off

REM Di chuyển vào thư mục oath-keeper
cd oath-keeper

REM Lệnh để build Docker image
docker build -t oath-keeper .

REM Lệnh để run Docker container ở chế độ background
docker run -d --rm --name authsvc-oathkeeper -p 4455:4455 -p 4456:4456 oath-keeper --config /config.yaml serve

REM Quay lại thư mục gốc (nếu cần thiết)
cd ..
