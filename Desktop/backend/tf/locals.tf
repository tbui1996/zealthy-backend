locals {
  users_internal_http_routes = {
    user_list            = "GET /users/user_list"
    revoke_access        = "PUT /users/revoke_access"
    get_organizations    = "GET /users/organizations"
    create_organizations = "POST /users/organizations"
    update_user          = "PUT /users/update_user"
  }

  users_external_http_routes = {
    sign_in = "POST /auth/external_sign_in"
    sign_up = "POST /auth/external_sign_up"
    refresh = "POST /auth/external_refresh"
  }

  support_internal_http_routes = {
    chat_session                = "POST /support/chat_session"
    pending_chat_sessions       = "GET /support/pending_chat_sessions"
    assign_pending_chat_session = "POST /support/assign_pending_chat_session"
    get_session_messages        = "GET /support/sessions/{id}/messages"
    get_sessions                = "GET /support/sessions"
    update_chat_session         = "POST /support/chat_session_update"
    update_chat_notes           = "PUT /support/{id}/notes"
  }

  support_external_http_routes = {
    send_message          = "POST /support/session/{id}/messages"
    get_session_messages  = "GET /support/session/{id}/messages"
    submit_feedback       = "POST /support/feedback"
    chat_session          = "POST /support/chat_session"
    get_user_sessions     = "GET /support/users/{id}/sessions"
    chat_session_star     = "PUT /support/session/{id}/star"
    online_internal_users = "GET /support/online"
    patients_get          = "GET /support/patients"
  }

  forms_internal_http_routes = {
    count    = "GET /forms/count"
    create   = "POST /forms"
    get      = "GET /forms/{id}"
    list     = "GET /forms"
    response = "GET /forms/{id}/response"
    send     = "POST /forms/{id}/send"
    delete   = "PUT /forms/{id}/delete"
    edit     = "PUT /forms/edit"
    close    = "PUT /forms/{id}/close"
  }

  router_internal_http_routes = {
    broadcast = "POST /router/broadcast"
  }

  deprecated_routes = {
    router_user_list = "GET /router/users"
  }

  cloud_internal_http_routes = {
    file_download         = "GET /cloud/file_download/{id}"
    file_upload           = "POST /cloud/upload"
    get_file              = "GET /cloud/get_file"
    associate_file        = "PUT /cloud/associate_file"
    delete_file           = "PUT /cloud/delete_file/{id}"
    pre_signed_upload_url = "GET /cloud/upload/url"
  }

  // TODO: DEPRECATED: Remove upload / download ASAP
  cloud_external_http_routes = {
    file_upload   = "POST /cloud/file_upload"
    file_download = "GET /cloud/file_download/{id}"
  }

  cloud_external_v2_http_routes = {
    file_upload           = "POST /cloud/upload"
    file_download         = "GET /cloud/file_download/{id}"
    pre_signed_upload_url = "GET /cloud/upload/url"
  }

  feature_flags_internal_http_routes = {
    create_flag = "POST /flags"
    list_flags  = "GET /flags"
    patch_flag  = "PATCH /flags/{id}"
    evaluate    = "GET /flags/evaluate"
    delete_flag = "DELETE /flags/{id}"
  }

  feature_flags_external_http_routes = {
    evaluate = "GET /flags/evaluate"
  }

  appointments_internal_http_routes = {
    list_appointments   = "GET /appointment"
    create_appointments = "POST /appointment"
    edit_appointments   = "PUT /appointment/{appointment_id}"
    delete_appointments = "DELETE /appointment/{appointment_id}"
  }

  patients_internal_http_routes = {
    list_patients   = "GET /patient"
    create_patients = "POST /patient"
    patch_patients  = "PUT /patient/{patient_id}"

  }

  agency_providers_internal_http_routes = {
    list_agency_providers   = "GET /agency_provider"
    create_agency_providers = "POST /agency_provider"
    edit_agency_providers   = "PUT /agency_provider/{agency_provider_id}"
  }

  live_envs_last_stable_commit = tomap({
    dev  = var.last_stable_develop_commit
    test = var.last_stable_test_commit
    prod = var.last_stable_production_commit
  })
}

locals {
  internal_groups = {
    program_manager = concat(
      // give access to internal users to the websocket api
      ["arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.webapp_websocket_api.id}/*/*"],
      // arn:aws:execute-api:region:account-id:api-id/stage/METHOD_HTTP_VERB/Resource-path
      [for route in concat(values(local.users_internal_http_routes), values(local.support_internal_http_routes), values(local.forms_internal_http_routes), values(local.router_internal_http_routes), values(local.cloud_internal_http_routes), values(local.appointments_internal_http_routes), values(local.agency_providers_internal_http_routes), values(local.patients_internal_http_routes)) : format("arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.gateway.id}/*/%s", replace(replace(route, "/\\s/", ""), "/\\{.*\\}/", "*"))]
    ),
    general_support = concat(
      ["arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.webapp_websocket_api.id}/*/*"],
      // general support does not need to send forms, broadcasts, or edit users
      [for route in concat(values(local.support_internal_http_routes), values(local.cloud_internal_http_routes), values(local.users_internal_http_routes)) : format("arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.gateway.id}/*/%s", replace(replace(route, "/\\s/", ""), "/\\{.*\\}/", "*"))]
    ),
    // might be more routes than we need in development_admin todo: clean them up
    development_admin = concat(
      // give access to internal users to the websocket api
      ["arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.webapp_websocket_api.id}/*/*"],
      // arn:aws:execute-api:region:account-id:api-id/stage/METHOD_HTTP_VERB/Resource-path
      [for route in concat(
        values(local.feature_flags_internal_http_routes),
        values(local.users_internal_http_routes),
        values(local.support_internal_http_routes),
        values(local.forms_internal_http_routes),
        values(local.router_internal_http_routes),
        values(local.cloud_internal_http_routes),
        values(local.appointments_internal_http_routes),
        values(local.patients_internal_http_routes),
        values(local.agency_providers_internal_http_routes)
        )
      : format("arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.gateway.id}/*/%s", replace(replace(route, "/\\s/", ""), "/\\{.*\\}/", "*"))]
    )
  }

  external_groups = {
    supervisor = concat(
      // give access to loop users to the websocket api
      ["arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.loop_websocket_api.id}/*/*"],
      // arn:aws:execute-api:region:account-id:api-id/stage/METHOD_HTTP_VERB/Resource-path
      // V1 Routes
      [for route in concat(
        values(local.users_external_http_routes),
        values(local.support_external_http_routes),
        values(local.cloud_external_http_routes),
        values(local.feature_flags_external_http_routes))
      : format("arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.loop_gateway.id}/*/%s", replace(replace(route, "/\\s/", ""), "/\\{.*\\}/", "*"))],
      // V2 Routes
      [for route in concat(
        values(local.cloud_external_v2_http_routes))
      : format("arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${module.loop_api_gateway_v2.api_gateway_id}/*/%s", replace(replace(route, "/\\s/", ""), "/\\{.*\\}/", "*"))]
    ),
    guest = concat(
      ["arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.loop_websocket_api.id}/*/*"],
      // V1 Routes
      [for route in concat(
        values(local.users_external_http_routes),
        values(local.support_external_http_routes),
        values(local.cloud_external_http_routes),
        values(local.feature_flags_external_http_routes))
      : format("arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${aws_apigatewayv2_api.loop_gateway.id}/*/%s", replace(replace(route, "/\\s/", ""), "/\\{.*\\}/", "*"))],
      // V2 Routes
      [for route in concat(
        values(local.cloud_external_v2_http_routes))
      : format("arn:aws:execute-api:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}:${module.loop_api_gateway_v2.api_gateway_id}/*/%s", replace(replace(route, "/\\s/", ""), "/\\{.*\\}/", "*"))]
    ),
  }
}
