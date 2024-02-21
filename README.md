# 중앙대학교 공지사항 RSS 서비스
개인적으로 쓰기 위해 만든 중앙대학교 공지사항 RSS 서비스입니다. 원래는 웹툰 RSS 기능도 있었지만 삭제했습니다. [golang](https://go.dev) 씁니다.

AWS Lambda 함수로 작동합니다. 빌드된 zip 파일로 AWS 람다 함수를 생성하세요.

## 빌드 / 테스트 방법
1. 버킷을 만들고 권한도 알아서 설정합니다.
1. 소스 코드 내에서 버킷 이름을 알아서 바꿉니다.
1. `./build-lambda-zip.sh`
1. 빌드된 zip 파일을 이용해서 람다를 생성하고 권한 잘 생성하면 끝!

## 환경변수
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