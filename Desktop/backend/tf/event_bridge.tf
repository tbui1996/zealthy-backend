resource "aws_cloudwatch_event_bus" "service_events" {
  name = "service-events"
}

### Start Schema Registry ###
## Typicall these schemas would be maintained in relative services, these schemas are common ##
resource "aws_schemas_registry" "sonar_common_schema_registry" {
  name        = "sonar_common_schema_registry"
  description = "Schema registry for events that come from more than one service."
}

resource "aws_schemas_schema" "connection_created_event" {
  name          = "connection_created_event"
  registry_name = aws_schemas_registry.sonar_common_schema_registry.name
  type          = "OpenApi3"
  description   = "Event fired when a user creates a new connection"

  content = jsonencode({
    "openapi" : "3.0.0",
    "info" : {
      "version" : "1.0.0",
      "title" : "Event"
    },
    "paths" : {},
    "components" : {
      "schemas" : {
        "Event" : {
          "type" : "object",
          "required" : ["UserID", "CreatedAt"],
          "properties" : {
            "UserID" : {
              "type" : "string",
            },
            "CreatedAt" : {
              "type" : "number",
              "format" : "epoch"
            },
          }
        }
      }
    }
  })
}
### End Schema Registry ###