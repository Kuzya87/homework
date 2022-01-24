terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
}

provider "yandex" {
  #token     = var.yc_token
  service_account_key_file = var.yc_service_account_key_file
  cloud_id  = var.yc_cloud
  folder_id = var.yc_folder
}