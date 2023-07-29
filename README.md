# 중앙대학교 공지사항 RSS 서비스
개인적으로 쓰기 위해 만든 중앙대학교 공지사항 RSS 서비스입니다. 원래는 웹툰 RSS 기능도 있었지만 삭제했습니다. [golang](https://go.dev) 씁니다.

## 빌드 / 테스트 방법
 - 빌드 : `go build`
 - 테스트 : `go test ./...` (참고: redis 기능 테스트하려면 redis 관련 환경변수 전달해줘야 함.)

## Redis 관련 환경변수
 - `REDIS_ENABLED`: redis를 이용하려면 `true`로 설정
 - `REDIS_ADDR`: redis 서버 주소 e.g. `127.0.0.1:6379`
 - `REDIS_DB` 이용할 redis DB 번호, 미지정시 0으로 지정됨.

## 서버 관련 환경변수
 - `PORT`: 서버포트
 - `GIN_MODE`: 프로덕션에서 돌리려면 `release`로 설정
 - `WEB_ADDRESS`: 웹사이트 주소 (예시: `https://rss.example.com`)

## 이미지 리소스 라이선스
### [html5-valid-badge](https://github.com/bradleytaunt/html5-valid-badge) License
```
MIT License

Copyright (c) 2019 Bradley Taunt

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## License
cau-rss --- RSS Service for Notices of Chuang-Ang University

Copyright (C) 2019~2023 Yeonjin Shin

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.