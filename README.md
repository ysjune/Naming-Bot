# Naming Bot

가끔 떠오르지 않는 변수명을 봇에게 떠맡기기 위해 만들었습니다.

## 사용 방법
```shell
@봇이름 .변수명 사용자 계정 상태
```

봇을 멘션한 후에 `.변수명` 이라는 커멘트 후에 변수명을 짓고 싶은 단어를 입력합니다.  


- 현재는 멘션을 지정해야 됩니다. 
  - 추후에는 멘션 안하고 커멘드만으로...
- 커멘드는 `.변수명` 만 지원하고 있습니다.
  - ~~다른 게 필요하면 직접 추가...하세요?~~

## 참고
이런 부분을 참고하여 작성하였습니다.
- [카카오페이 너의 이름은 변수명 챗봇 개발](https://tech.kakaopay.com/post/variable-name-bot-haero-sery-bread/)
- [go-openai](https://github.com/sashabaranov/go-openai)
- [slack-go](https://github.com/slack-go/slack)
- [slack bot using golang](https://www.bacancytechnology.com/blog/develop-slack-bot-using-golang)