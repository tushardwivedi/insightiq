import os

# Database configuration - Build URI from environment variables
POSTGRES_USER = os.environ.get('POSTGRES_USER', 'superset')
POSTGRES_PASSWORD = os.environ.get('POSTGRES_PASSWORD', 'superset')
POSTGRES_DB = os.environ.get('POSTGRES_DB', 'superset')
SQLALCHEMY_DATABASE_URI = f'postgresql://{POSTGRES_USER}:{POSTGRES_PASSWORD}@postgres:5432/{POSTGRES_DB}'

# Security
SECRET_KEY = os.environ.get('SECRET_KEY', 'insightiq-secret-key-change-in-production')
WTF_CSRF_ENABLED = False

# Enable features
FEATURE_FLAGS = {
    "EMBEDDED_SUPERSET": True,
    "ENABLE_TEMPLATE_PROCESSING": True,
}

# CORS settings
ENABLE_CORS = True
CORS_OPTIONS = {
    'supports_credentials': True,
    'allow_headers': ['*'],
    'resources': {'*': {'origins': '*'}},
}

# API access
GUEST_ROLE_NAME = 'Gamma'
PUBLIC_ROLE_LIKE_GAMMA = True