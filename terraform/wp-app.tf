locals {
  wp-app-vm-list = tolist(concat(yandex_compute_instance.wp-app-a[*],
    yandex_compute_instance.wp-app-b[*],
    yandex_compute_instance.wp-app-c[*]))
}

resource "yandex_compute_instance" "wp-app-a" {
  count = var.count_compute_instances
  name = "wp-app-a-${count.index}"
  zone = "ru-central1-a"

  resources {
    cores  = 2
    memory = 2
  }

  boot_disk {
    initialize_params {
      image_id = "fd80viupr3qjr5g6g9du"
    }
  }

  network_interface {
    subnet_id = yandex_vpc_subnet.wp-subnet-a.id
    nat       = true
  }

  metadata = {
    ssh-keys = "ubuntu:${file("~/.ssh/yc.pub")}"
  }
}

resource "yandex_compute_instance" "wp-app-b" {
  count = var.count_compute_instances
  name = "wp-app-b-${count.index}"
  zone = "ru-central1-b"

  resources {
    cores  = 2
    memory = 2
  }

  boot_disk {
    initialize_params {
      image_id = "fd80viupr3qjr5g6g9du"
    }
  }

  network_interface {
    subnet_id = yandex_vpc_subnet.wp-subnet-b.id
    nat       = true
  }

  metadata = {
    ssh-keys = "ubuntu:${file("~/.ssh/yc.pub")}"
  }
}

resource "yandex_compute_instance" "wp-app-c" {
  count = var.count_compute_instances
  name = "wp-app-c-${count.index}"
  zone = "ru-central1-c"

  resources {
    cores  = 2
    memory = 2
  }

  boot_disk {
    initialize_params {
      image_id = "fd80viupr3qjr5g6g9du"
    }
  }

  network_interface {
    subnet_id = yandex_vpc_subnet.wp-subnet-c.id
    nat       = true
  }

  metadata = {
    ssh-keys = "ubuntu:${file("~/.ssh/yc.pub")}"
  }
}