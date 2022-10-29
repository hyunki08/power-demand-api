# Power Demand Chart api

2007년 ~ 2021년 제주지역 일별 시간단위 전력 수요량 차트 앱의 rest api

### 구현 스택

- Golang
- Gin
- MongoDB
- Docker / Docker Compose

### 구현 내용

db, api 서비스를 `dockerfile`로 빌드 설정 후 `docker-compose` 로 실행 할 수 있도록함

**db 구축**

- Docker에 mongo 컨테이너 생성
- db 컨테이너를 로컬 볼륨과 연결하기 보다는 db실행 시 csv를 바로 바로 적재하도록 함
- 전력 수요량 csv를 db에 추가하기 위해 스크립트 작성
  - seed.sh
    ```bash
    #!/bin/bash
    mongoimport --db PDDB --collection PowerDemand --type csv --headerline --file /files/power-demand.csv
    ```
- Dockerfile 에 csv, seed.sh 파일을 복사하도록 추가
  - seed.sh는 /docker-entrypoint-initdb.d에 복사하여 db가 실행 된 후 스크립트가 실행되도록 함
    [mongo - Official Image | Docker Hub](https://hub.docker.com/_/mongo)
    **Initializing a fresh instance 항목 확인**
  - DB : PDDB
  - Collection : PowerDemand

**API 서버 구축**

- Go의 gin 라이브러리를 통해 서버 구현
  - air 패키지를 통해 hot-reload 처리
- mongodb driver 패키지를 설치하여 db정보 불러옴
  - godotenv 패키지를 통해 `.env` 처리

### Frontend

https://github.com/hyunki08/power-demand-app
