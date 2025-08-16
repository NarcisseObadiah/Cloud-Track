# Cloud PaaS Platform

A comprehensive Platform-as-a-Service (PaaS) solution built on OpenStack infrastructure with Kubernetes orchestration, providing managed database services with automated provisioning, monitoring, and multi-tenant support.

![Platform Architecture](/paas.png)

## 🚀 Project Overview


## 🖥️ Frontend User Interface

The Vue.js frontend provides a comprehensive web interface for managing the PaaS platform:

### System Administrator Dashboard

**Features:**
- Real-time cluster health monitoring
- Resource utilization overview
- Active tenant management
  
![System Admin Dashboard](images/sys-admin.png)

### Database Management Interface
![Database Creation Form](/images/db-mngt.png)

### Database instance form
- Interactive database creation wizard
![Database Creation Form](/images/db-creaton-on.png)

**Response after Database Creation Request**
- Connection string generation
- Hashed password provided
- Multi environnement config provided to authenticate to the database instance
![Database Creation Form](/images/db-response.png)

### Tenant Management Console
**User Features:**
- Personal database dashboard
- Resource usage tracking
- Self-service database operations
![Tenant Console](/images/Tenants.png)

### Pod and Service Monitoring
**Monitoring Views:**
- Real-time pod status grid
- Service health indicators
- Log streaming interface
- Resource consumption charts
  
![Pod Status View](/images/admin-pod-overview.png)


### User Authentication Third Party through Zitadel

**Security Features:**
- Zitadel SSO integration
- Role-based access control
- Session management
- Multi-factor authentication support
  
![Login Interface](/images/add-user.png)


This project implements a complete cloud platform following a three-phase development approach:

### Phase 1: IaaS Foundation with OpenStack
- **OpenStack Environment**: Complete setup with Nova, Neutron, Glance, and Keystone services
- **Infrastructure Management**: Automated VM, network, and storage provisioning
- **Identity & Access Management**: Zitadel integration for RBAC and authentication

### Phase 2: PaaS Layer with Kubernetes
- **Kubernetes Cluster**: Deployed on OpenStack infrastructure using kubeadm
- **Managed Database Service**: PostgreSQL-as-a-Service with Zalando operator
- **RESTful API**: Go-based backend for service provisioning and management
- **Resource Metering**: CPU, memory, and storage consumption tracking

### Phase 3: Advanced Platform Features
- **Auto-scaling**: Horizontal Pod Autoscaler (HPA) implementation
- **Multi-region Deployment**: High availability and disaster recovery
- **Web UI**: Vue.js frontend with JWT authentication
- **Observability**: Prometheus, and  Grafana integration
- **Audit Logging**: Comprehensive security and compliance tracking (must have)

![Admin Dashboard Overview](/images/welcome.png)

## 🏗️ Architecture
### Infrastructure Components
-**Terraforn For Provisioning**
![Infrastructure Diagram](/images/terra-for-infra.png)

- **OpenStack Cluster**: 3-node setup (1 master, 2 workers)
- **Kubernetes Cluster**: K3s distribution for lightweight orchestration (Ansible)
 ![Infrastructure Diagram](/images/full-infra.png)

- **Database Service**: Zalando PostgreSQL Operator for managed databases
  ![Database Management UI](/images/db-cluster-overview.png)
  
- Pod Management Dasboard, after creation
![Pod Management Dashboard](/images/Pods_Mgmt.png)

### Technology Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| **Infrastructure** | OpenStack (Nova, Neutron, Glance, Keystone) | IaaS foundation |
| **Orchestration** | Kubernetes (K3s) | Container orchestration |
| **Database** | PostgreSQL + Zalando Operator | Managed database service |
| **Backend** | Go + Gin Framework | REST API development |
| **Frontend** | Vue.js | User interface |
| **Monitoring** | Prometheus + Grafana | Metrics and visualization |
| **Logging** | Loki | Log aggregation |
| **Authentication** | Zitadel + JWT | Identity management |
| **IaC** | Terraform + Ansible | Infrastructure automation |

## 🛠️ Installation & Setup

### Prerequisites

- OpenStack environment with admin access
- Terraform >= 1.0
- Ansible >= 2.9
- kubectl
- Docker
- Go >= 1.19

### 1. Infrastructure Provisioning

```bash
# Clone the repository
git clone <repository-url>
cd paas-api

# Configure Terraform variables
cp terraform/terraform.tfvars.example terraform/terraform.tfvars

# Deploy infrastructure
cd terraform
terraform init
terraform plan
terraform apply
```

### 2. Kubernetes Cluster Setup

```bash
# Configure Ansible inventory
cp ansible/hosts.ini.example ansible/hosts.ini
# Update hosts.ini with my node IPs

# Install Kubernetes cluster
cd ansible
ansible-playbook -i hosts.ini install-k3s.yml
```

### 3. Deploy PaaS Services

```bash
# Install Zalando PostgreSQL Operator
kubectl apply -k https://github.com/zalando/postgres-operator/manifests

# Deploy monitoring stack
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install monitoring prometheus-community/kube-prometheus-stack \
  --namespace monitoring --create-namespace

# Build and deploy API
docker build -t paas-api:latest .
kubectl apply -f k8s/deployment.yaml
```

## 📚 API Documentation

### Database Management Endpoints

#### Create Database
```http
POST /databases
Content-Type: application/json

{
  "username": "testuser",
  "db_name": "myapp",
  "replicas": 1
}
```

**Response:**
```json
{
  "database_name": "myapp",
  "host": "myapp.tenant-testuser.svc.cluster.local",
  "port": "5432",
  "status": "provisioning",
  "message": "Database is being created. Credentials will be available shortly.",
  "secret_name": "myapp.myapp.credentials.postgresql.acid.zalan.do",
  "instructions": {
    "check_status": "kubectl get postgresql myapp -n tenant-testuser",
    "get_credentials": "kubectl get secret myapp.myapp.credentials.postgresql.acid.zalan.do -n tenant-testuser -o yaml"
  }
}
```

#### List Databases
```http
GET /databases/{namespace}
```

#### Get Database Status
```http
GET /databases/{namespace}/{db_name}/status
```

#### Delete Database
```http
DELETE /databases/{namespace}/{db_name}
```

### Infrastructure Endpoints

#### List Pods
```http
GET /pods
GET /pods/{namespace}
```


## 🔧 Configuration

### Environment Variables

```bash
# Kubernetes Configuration
export KUBECONFIG=/path/to/kubeconfig

# API Configuration
export PORT=8080
export GIN_MODE=release

# OpenStack Credentials
export OS_AUTH_URL=http://your-openstack:5000/v3
export OS_USERNAME=admin
export OS_PASSWORD=your-password
export OS_PROJECT_NAME=admin
export OS_REGION_NAME=RegionOne
```

### Database Template Configuration

The PostgreSQL clusters are configured via templates in `templates/postgres-cluster.yaml.tmpl`:

```yaml
apiVersion: "acid.zalan.do/v1"
kind: postgresql
metadata:
  name: {{ .DBName }}
  namespace: {{ .Namespace }}
spec:
  teamId: "{{ .Team }}"
  volume:
    size: 5Gi
  numberOfInstances: {{ .Replicas }}
  resources:
    requests:
      cpu: 200m
      memory: 256Mi
    limits:
      cpu: 500m
      memory: 512Mi
```

## 📊 Monitoring & Observability

### Prometheus Metrics

Access Prometheus at: `http://localhost:9090` (with port-forward)

```bash
kubectl port-forward svc/monitoring-kube-prometheus-prometheus -n monitoring 9090:9090
```

### Grafana Dashboards

Access Grafana at: `http://localhost:3000` (with port-forward)

```bash
kubectl port-forward svc/monitoring-grafana -n monitoring 3000:80
```

**Default Credentials:**
- Username: `admin`
- Password: `prom-operator`

![Grafana Dashboard](/images/Grafana-dashboard.png)

### Key Metrics Monitored

- **Golden Metrics**: Latency, Traffic, Errors, Saturation
- **Infrastructure**: Node CPU, Memory, Disk usage
- **Application**: API response times, database connections
- **Kubernetes**: Pod health, resource utilization

## 🔐 Security & Authentication

### Zitadel Integration

1. **User Authentication**: Login via Zitadel IdP
2. **JWT Token**: Secure token generation
3. **RBAC**: Role-based access control
4. **API Authorization**: Bearer token validation

![Authentication Flow](/images/Zitadel_Conf.png)

### Security Groups
OpenStack security groups configured for:
- **SSH Access**: Port 22 for management
- **API Access**: Port 8080 for PaaS API
- **Kubernetes**: Ports 6443, 10250 for cluster communication
- **Database**: Port 5432 for PostgreSQL (internal only)

![Infrastructure Diagram](/images/security_group.png)

### Namespace Isolation

Each tenant gets a dedicated Kubernetes namespace:
```
tenant-{username}/
├── PostgreSQL instances
├── Secrets (credentials)
├── ConfigMaps
└── NetworkPolicies
```

### Resource Quotas

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: tenant-quota
spec:
  hard:
    requests.cpu: "2"
    requests.memory: 4Gi
    persistentvolumeclaims: "5"
    pods: "10"
```

## 🚀 Deployment Options

### Local Development

```bash
# Run API locally
go run main.go

# Access via port-forward
kubectl port-forward deployment/paas-api 8080:8080
```

### Production Deployment

```bash
# Build and push Docker image
docker build -t narcisse198/paas-api:v1.0.0 .
docker push narcisse198/paas-api:v1.0.0

# Deploy via Kubernetes
kubectl apply -f k8s/production/
```

### NodePort Access

Access the API externally via NodePort:
```bash
# Find node IPs
kubectl get nodes -o wide

# Access API
curl -X POST http://<NODE-IP>:30971/databases \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","db_name":"mydb"}'
```

## 📈 Cluster Monitoring

![Resource Monitoring Dashboard](/images/cluster-monitoring.png)

### Horizontal Pod Autoscaler (ideal implementation but not provided in the current repos)

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: paas-api-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: paas-api
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

### Database Scaling (Ideal Options)

PostgreSQL instances support:
- **Vertical Scaling**: CPU/Memory resource adjustments
- **Horizontal Scaling**: Read replicas (up to 5 instances)
- **Storage Scaling**: Dynamic PVC expansion

## 🛡️ Backup & Disaster Recovery

### Database Backups

```bash
# Automated backups via CronJob
kubectl apply -f k8s/backup/postgres-backup-cronjob.yaml
```

### Multi-Region Setup (Not implemented in the current repos but has been implemented for the final Group Presentation At StackIT)

- **Primary Region**: Full stack deployment
- **DR Region**: Standby with data replication
- **Failover**: Automated traffic routing

## 📝 Troubleshooting

### Common Issues

#### API Server Timeouts
```bash
# Check API server health
kubectl cluster-info

# Verify SSH tunnel (if using remote cluster)
ssh -L 6443:localhost:6443 user@master-node
```

For support and questions:
- 📧 Email: narcisseobadiahdm@gmail.com
---

**Built with ❤️ by NarcisseObadiah using OpenStack, Kubernetes, and Go**
