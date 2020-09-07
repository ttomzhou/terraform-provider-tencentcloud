provider "tencentcloud" {
  region = "ap-guangzhou"
}

#resource "tencentcloud_image" "image_snap" {
#   image_name          = var.image_snapshot_name
#   snapshot_ids        = ["snap-nbp3xy1d", "snap-nvzu3dmh"]
#   force_power_off     = true
#   image_description   = "create image with snapshot"
#}

resource "tencentcloud_image" "image_instance" {
   image_name         = var.image_imstance_name
   instance_id        = "ins-2ju245xg"
   data_disk_ids      = ["disk-gii0vtwi"]
   force_power_off    = true
   sysprep            = false
   image_description  = "create image with instance"
}
