resource "aws_ses_template" "new_message_open" {
  name    = "sonar_new_message_open"
  subject = "You have a new message!"
  html    = file("../email-templates/new_message_open.html")
}

resource "aws_ses_template" "feedback" {
  name    = "sonar_feedback"
  subject = "Greetings, {{name}}"
  html    = file("../email-templates/sonar_feedback.html")
}