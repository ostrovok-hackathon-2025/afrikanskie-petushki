CREATE MATERIALIZED VIEW IF NOT EXISTS monthly_stats AS
SELECT 
    DATE_TRUNC('month', CURRENT_DATE) as month_start,
    
    COUNT(DISTINCT CASE WHEN o.status = 'done' 
          AND o.check_in_at >= DATE_TRUNC('month', CURRENT_DATE)
          AND o.check_in_at < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month'
          THEN o.id END) as completed_offers,
    
    COUNT(DISTINCT CASE WHEN a.created_at >= DATE_TRUNC('month', CURRENT_DATE)
          AND a.created_at < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month'
          THEN a.id END) as applications_received,
    
    COUNT(DISTINCT CASE WHEN r.status = 'accepted'
          AND r.created_at >= DATE_TRUNC('month', CURRENT_DATE)
          AND r.created_at < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month'
          THEN r.id END) as accepted_reports

FROM offer o
LEFT JOIN application a ON o.id = a.offer_id
LEFT JOIN report r ON a.id = r.application_id;