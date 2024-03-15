pipeline {
  agent any
  stages {
    stage('Send Slack Message') {
      steps {
        slackSend(baseUrl: 'https://jenkins.dev.spandigital.io/', color: '#5cc9f5', message: 'Test Message', attachments: 'nothing', blocks: 'nothing')
      }
    }

  }
}