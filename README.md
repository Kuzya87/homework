# Описание выполненного домашнего задания

1. По основной части ДЗ всё было успешно выполнено:

        Apply complete! Resources: 9 added, 0 changed, 0 destroyed.

        Outputs:

        database_host_fqdn = tolist([
          "rc1b-u6n53aehka3erolg.mdb.yandexcloud.net",
          "rc1c-13fpnl73qqd5t9hg.mdb.yandexcloud.net",
        ])
        load_balancer_public_ip = tolist([
          "51.250.27.61",
        ])
    
    Удаление также прошло успешно.

2. Задание со звёздочкой успешно выполнено - вручную создан сервисный аккаунт, ему назначены права на каталог и Terraform всё выполняет от его имени.

3. Задание с двумя звёдочками выполнено.

    Удалось использовать единственный динамический блок target для yandex_lb_target_group, так как я загнал в одну local-переменную wp-app-vm-list список цельных объектов всех ВМ wp-app и прошёлся по этому списку через for_each.