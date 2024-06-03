@echo off

REM Di chuyển vào thư mục oath-keeper
cd oath-keeper

REM Lệnh để build Docker image
docker stop authsvc-oathkeeper

REM Quay lại thư mục gốc (nếu cần thiết)
cd ..
