import { Container, Card, Link, Row, Col, Text, Spacer, Input, Modal, Button, styled, Loading } from "@nextui-org/react";
import { useState } from 'react';

export const SendButton = styled('button', {
  // reset button styles
  background: 'transparent',
  border: 'none',
  padding: 0,
  margin: 0,
  // styles
  width: '24px',
  margin: '0 10px',
  dflex: 'center',
  bg: '$primary',
  borderRadius: '$rounded',
  cursor: 'pointer',
  transition: 'opacity 0.25s ease 0s, transform 0.25s ease 0s',
  svg: {
    size: '100%',
    padding: '4px',
    transition: 'transform 0.25s ease 0s, opacity 200ms ease-in-out 50ms',
    boxShadow: '0 5px 20px -5px rgba(0, 0, 0, 0.1)',
  },
  '&:hover': {
    opacity: 0.8
  },
  '&:active': {
    transform: 'scale(0.9)',
    svg: {
      transform: 'translate(24px, -24px)',
      opacity: 0
    }
  }
});

export const SendIcon = ({
  fill = "currentColor",
  filled,
  size,
  height,
  width,
  label,
  className,
  ...props
}) => {
  return (
    <svg
      data-name="Iconly/Curved/Lock"
      xmlns="http://www.w3.org/2000/svg"
      width={ size || width || 24 }
      height={ size || height || 24 }
      viewBox="0 0 24 24"
      className={ className }
      { ...props }
    >
      <g transform="translate(2 2)">
        <path
          d="M19.435.582A1.933,1.933,0,0,0,17.5.079L1.408,4.76A1.919,1.919,0,0,0,.024,6.281a2.253,2.253,0,0,0,1,2.1L6.06,11.477a1.3,1.3,0,0,0,1.61-.193l5.763-5.8a.734.734,0,0,1,1.06,0,.763.763,0,0,1,0,1.067l-5.773,5.8a1.324,1.324,0,0,0-.193,1.619L11.6,19.054A1.91,1.91,0,0,0,13.263,20a2.078,2.078,0,0,0,.25-.01A1.95,1.95,0,0,0,15.144,18.6L19.916,2.525a1.964,1.964,0,0,0-.48-1.943"
          fill={ fill }
        />
      </g>
    </svg>
  );
};

export default function App() {
  const [visible, setVisible] = useState(false);
  const handler = () => setVisible(true);
  const closeHandler = () => {
    setVisible(false);
    console.log("closed");
  };
  const [values, setValues] = useState({
    twitter_id: '',
    content: '',
  })
  const [loading, setLoading] = useState(false);
  const handleChange = (prop) => (event) => {
    setValues({ ...values, [prop]: event.target.value })
  }

  const handleConfirm = async () => {
    setVisible(false);
    setLoading(true)
    try {
      const response = await fetch('https://tweet-api-boe.aireview.tech/api/get_tweet_analysis?twitter_id=' + values.twitter_id, {
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'GET',
      })
      if (!response.ok) {
        const error = await response.json()
        console.error(error.error)
        setCurrentError(error.error)
        throw new Error('Request failed')
      }
      const data = response.body
      if (!data)
        throw new Error('No data')

      const reader = data.getReader()
      const decoder = new TextDecoder('utf-8')
      let done = false

      while (!done) {
        const { value, done: readerDone } = await reader.read()
        if (value) {
          const char = decoder.decode(value)

          if (char.startsWith('Data:')) {
            setValues({ content: char.substring(5) })
          } else {
            setValues({ content: values.content + char })
          }

          // const lines = char.trim().split('\n')
          // console.log(lines)

          // var isData = false

          // for (let i = 0; i < lines.length; i++) {
          //   if (lines[i] == 'Data:') {
          //     isData = true
          //     setValues({ content: '' })
          //     continue
          //   }

          //   if (isData == true) {
          //     const line = lines[i]
          //     setValues({ content: values.content + line })
          //     i++
          //     continue
          //   }

          //   if (isData == false) {
          //     const line = lines[i]
          //     setValues({ content: values.content + line })
          //     continue
          //   }
          // }

          // for (let i = 0; i < lines.length; i++) {
          //   if (lines[i] === '') {
          //     i++
          //     continue
          //   }
          //   if (lines[i].startsWith('event:plugin_selected')) {
          //     const line = lines[i + 1]
          //     const data = line.substring(5)
          //     // setCurrentPlugin(data)
          //     const plugin = `Plugin: \`\`\`${data}\`\`\`` + '\n'
          //     setCurrentAssistantMessage(currentAssistantMessage() + plugin)
          //     i++
          //     continue
          //   }
          //   if (lines[i].startsWith('event:plugin_requested')) {
          //     const line = lines[i + 1]
          //     const data = line.substring(5)
          //     // setCurrentPluginReq(data)
          //     const splitData = data.split(/(.{50})/).filter(Boolean)
          //     console.log('splitData', splitData)
          //     const req = `Plugin: \`\`\`json${splitData.join('\n')}\`\`\`` + '\n'
          //     // const req = `Plugin Req: \`\`\`${data}\`\`\`` + '\n'
          //     setCurrentAssistantMessage(currentAssistantMessage() + req)
          //     i++
          //     continue
          //   }
          //   if (lines[i].startsWith('event:plugin_responsed')) {
          //     const line = lines[i + 1]
          //     const data = line.substring(5)
          //     // setCurrentPluginResp(data)
          //     const resp = `Plugin Resp: \`\`\`${data}\`\`\`` + '\n'
          //     setCurrentAssistantMessage(currentAssistantMessage() + resp)

          //     i++
          //     continue
          //   }
          //   if (lines[i].startsWith('event:completion')) {
          //     const line = lines[i + 1]
          //     const data = JSON.parse(line.substring(5))
          //     setCurrentAssistantMessage(currentAssistantMessage() + data.delta)
          //     i++
          //     continue
          //   }
          // }
        }
        done = readerDone
        setLoading(false)
      }
    } catch (e) {
      console.error(e)
      return
    }
  };


  return (
    <Container sm display="flex" gap={7} 
      css={{
        marginTop: '4em'
      }}
    >
      <Col css={{
        textAlign: 'center'
      }}>
        <Text
          transform="full-size-kana"
          h1
          size={50}
          css={ {
            textGradient: "to right, #006E3A 8%, #166BB5 100%",
            '@xsMax': {fontSize: '2.5rem'},
          } }
          weight="bold"
        >
          推文分析器
        </Text>
        <Text
          transform="full-size-kana"
          h1
          size={ 37 }
          css={ {
            textGradient: "to right, #006E3A 8%, #166BB5 100%",
          } }
          weight="bold"
        >
          Tweet Analyzer
        </Text>
      </Col>
      <Spacer y={4} />
      <Card gap={ 2 }>
        <Card.Body>
          <Col justify="center" align="center">
            <Text css={{
              color: '$accents7'
            }}>请输入 Twitter ID [如 https://twitter.com/jack 即应输入 jack]</Text>
            <Spacer y={0.5}/>
            <Input
              clearable
              contentRightStyling={ false }
              label=""
              placeholder="L_x_x_x_x_x"
              labelLeft="ID"
              onChange={ handleChange('twitter_id') }
              value={ values.twitter_id }
              contentRight={
                !loading ?
                <SendButton onClick={ handler }>
                  <SendIcon />
                </SendButton>
                :
                <Loading size="sm" css={{margin: '.5em'}} />
              }
            />
            {
              loading && 
              <Text css={{
                marginTop: '$4',
                color: '$blue800',
                fontWeight: 'lighter'
              }}>由于当前请求人数较多，可能需要等待一段时间~</Text>
            }
          </Col>
          <Modal
            closeButton
            animated={ false }
            aria-labelledby="modal-title"
            open={ visible }
            onClose={ closeHandler }
          >
            <Modal.Header>
              <Text id="modal-title" size={ 18 }>
                QaQ 可以请你帮我点一个&nbsp;
                <Text b size={ 18 }>
                  关注
                </Text>
                &nbsp;耶
              </Text>
            </Modal.Header>
            <Modal.Body>
              <Text>你只需要点一下</Text>
              <Link color href="https://twitter.com/L_x_x_x_x_x" target="_blank">
                这里
              </Link>
              <Text>就好啦，非常感谢~</Text>
            </Modal.Body>
            <Modal.Footer>
              <Button bordered color="gradient" auto onPress={ handleConfirm }>
                点这里继续～
              </Button>
            </Modal.Footer>
          </Modal>
          <Spacer y={ 1 } />
          <Row justify="center" align="center">
            <Text h6 size={ 15 } color="black" css={ { m: 2 } }>
              { values.content }
            </Text>
          </Row>
        </Card.Body>
      </Card>
    </Container>
  );
}
