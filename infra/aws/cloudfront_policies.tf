# Política de cache que desativa o cache para rotas da API
resource "aws_cloudfront_cache_policy" "no_cache_api" {
  name = "no-cache-api-policy"

  default_ttl = 0
  max_ttl     = 0
  min_ttl     = 0

  parameters_in_cache_key_and_forwarded_to_origin {
    cookies_config {
      cookie_behavior = "none"
    }

    headers_config {
      header_behavior = "none"
    }

    query_strings_config {
      query_string_behavior = "none"
    }
  }
}

# Política de request para encaminhar tudo ao ALB
resource "aws_cloudfront_origin_request_policy" "api_origin_request" {
  name = "api-origin-request-policy"

  cookies_config {
    cookie_behavior = "all"
  }

  headers_config {
    header_behavior = "whitelist"
    headers {
      items = ["Host"]
    }
  }

  query_strings_config {
    query_string_behavior = "all"
  }
}
