##In case of errors...
Known errors all revolve around not including parallelism=1 when updating header links or redactions.  If you accidentally did this here are the known scenarios and how to get out.

##### My resource didn't actually update
```
Running subsequent terraform applys should fix the issue.
```

##### I get "Error: Not Found" when running plan or apply
```
Find out which resource is causing the issue and ensure that the resource does not exist in the console. 
Then you must manually edit the terraform.tfstate file and remove the specific resource from the list. 
You should then be able to proceed with a terraform plan/apply.
```

##### Terraform didn't delete my resource
```
You must manually delete these resources in the console. Unfortunately terraform does not always error out here.
```
