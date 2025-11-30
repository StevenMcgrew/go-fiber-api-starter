INSERT INTO notifications (text_content, has_viewed, user_id)
SELECT
    -- Generate sample text content
    'Notification ' || generate_series,
    -- Randomly assign TRUE or FALSE for HasViewed (approx. 50/50 split)
    CASE WHEN (generate_series % 2) = 0 THEN TRUE ELSE FALSE END,
    -- Assign a static or random UserId (here using static 1)
    1
FROM
    generate_series(1, 100);
