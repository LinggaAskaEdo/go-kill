db.user_activities.createIndex(
    { "user_id": 1, "timestamp": -1 },
    { name: "idx_user_activities_user_id_timestamp", background: true }
);
