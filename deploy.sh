
kubectl apply -f namespace-dev.yaml
kubectl apply -f zookeeper-deployment.yaml
kubectl apply -f zookeeper-service.yaml
kubectl apply -f kafka-deployment.yaml
kubectl apply -f kafka-service.yaml
kubectl apply -f ./product-service/product-service-deployment.yaml
kubectl apply -f ./product-service/product-service-service.yaml
kubectl apply -f ./order-service/order-service-deployment.yaml
kubectl apply -f ./order-service/order-service-service.yaml
kubectl apply -f ./notification-service/notification-service-deployment.yaml
kubectl apply -f ./notification-service/notification-service-service.yaml
