# Utiliser une image de base légère, par exemple alpine
FROM alpine:latest

# Définir un répertoire de travail dans le conteneur
WORKDIR /app

# Copier le binaire et le fichier de configuration dans l'image Docker
COPY bin/main /app/main
COPY bin/config/config.yml /app/config/config.yml

# Assurez-vous que le binaire a les permissions d'exécution
RUN chmod +x /app/main

# Passer Gin en mode production
export GIN_MODE=release

# Définir la commande par défaut pour exécuter le binaire
CMD ["/app/main"]
