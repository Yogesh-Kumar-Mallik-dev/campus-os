import { defineConfig, env } from 'prisma/config';

export default defineConfig({
  schema: 'prisma/schema',
  datasource: {
    url: env('DATABASE_URL') || 'postgresql://campus_admin:campus_password_dev@localhost:5434/campus_os?schema=auth_schema',
  },
});
