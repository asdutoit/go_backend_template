pipeline {
  agent any
  stages {
    stage('Send Slack Message') {
      steps {
        script {
          def attachments = [
            [
              text: 'This is a test message',
              fallback: 'Fallback message',
              color: '#ff0000'
            ]
          ]

          slackSend(channel: "#span-devops-feed", attachments: attachments)
        }
      }
    }
  }
}
