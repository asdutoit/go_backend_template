pipeline {
  agent any
  stages {
    stage('Send Slack Message') {
      steps {
        def attachments = [
          [
            text: 'I find your lack of faith disturbing!',
            fallback: 'Hey, Vader seems to be mad at you.',
            color: '#ff0000'
          ]
        ]

        slackSend(channel: "#span-devops-feed", attachments: attachments)
      }
    }

  }
}
