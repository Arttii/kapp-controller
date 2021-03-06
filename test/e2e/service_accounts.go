package e2e

import (
  "fmt"
)

type ServiceAccounts struct {
  namespace string
}

func (sa ServiceAccounts) ForNamespaceYAML() string {
  return fmt.Sprintf(`
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kappctrl-e2e-ns-sa
  annotations:
    kapp.k14s.io/change-rule: "delete after deleting kappctrl-e2e.k14s.io/apps"
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kappctrl-e2e-ns-role
  annotations:
    kapp.k14s.io/change-rule: "delete after deleting kappctrl-e2e.k14s.io/apps"
rules:
- apiGroups: ["*"]
  resources: ["*"]
  verbs: ["*"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: kappctrl-e2e-ns-role-binding
  annotations:
    kapp.k14s.io/change-rule: "delete after deleting kappctrl-e2e.k14s.io/apps"
subjects:
- kind: ServiceAccount
  name: kappctrl-e2e-ns-sa
  namespace: %s
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: kappctrl-e2e-ns-role
`, sa.namespace)
}
