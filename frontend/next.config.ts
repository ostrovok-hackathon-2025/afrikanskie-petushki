import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  /* config options here */
  experimental: {
    optimizeCss: false, // Отключаем lightningcss
  },
  
  output: 'standalone',
};

export default nextConfig;
