locals {
  app = "lofyd"

  function_name = "lofyd-forgotpasswordfunc"

  defaults ={
    accountId = "<set accountId in env_map>"
  }

  env_map = {
    dev={
      accountId = 277214047248
    }
    prod={

    }
  }
  configuration             = lookup(local.env_map, "${var.stage}", {})
  accountId                 = lookup(local.configuration, "accountId", local.defaults.accountId)
}