# Message Flow

## Table of Contents

- [Context](#context)
- [Services](#services)
  - [Analytics Service](#analytics-service)
  - [Campaign Service](#campaign-service)
  - [Notification Service](#notification-service)
  - [User Service](#user-service)
- [Channels](#channels)
  - [analytics.alert](#analyticsalert)
  - [analytics.insights](#analyticsinsights)
  - [analytics.report.request](#analyticsreportrequest)
  - [campaign.analytics](#campaignanalytics)
  - [campaign.create](#campaigncreate)
  - [campaign.execute](#campaignexecute)
  - [campaign.status](#campaignstatus)
  - [notification.analytics](#notificationanalytics)
  - [notification.preferences.get](#notificationpreferencesget)
  - [notification.preferences.update](#notificationpreferencesupdate)
  - [notification.user.{user_id}.push](#notificationuseruser_idpush)
  - [user.analytics](#useranalytics)
  - [user.info.request](#userinforequest)
  - [user.info.update](#userinfoupdate)

## Context
```mermaid
flowchart TD
    Analytics_Service["<b>Analytics Service</b><hr/>A centralized analytics service that receives and processes analytics events from all other services.
Provides insights, reporting, and analytics data aggregation for user behavior, notification performance,
campaign effectiveness, and system-wide metrics.
"]
    Campaign_Service["<b>Campaign Service</b><hr/>A service that manages notification campaigns, user targeting, and campaign execution.
Handles campaign creation, user segmentation, scheduling, and personalized notification delivery.
Uses user data for targeting and personalization of campaign messages.
"]
    Notification_Service["<b>Notification Service</b><hr/>A service that handles user notifications, preferences, and interactions.
Supports real-time notifications, user preferences management.
"]
    User_Service["<b>User Service</b><hr/>A service that manages user information, profiles, and authentication.
Handles user data requests, profile updates, and user lifecycle events.
"]
    Campaign_Service -->|"Pub"| Analytics_Service
    Campaign_Service -->|"Pub"| Notification_Service
    Campaign_Service -->|"Req"| User_Service
    Notification_Service -->|"Pub"| Analytics_Service
    Notification_Service <-->|"Pub/Req"| User_Service
    User_Service -->|"Pub"| Analytics_Service
    style Analytics_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Campaign_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Notification_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style User_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

## Services

### Analytics Service

A centralized analytics service that receives and processes analytics events from all other services.
Provides insights, reporting, and analytics data aggregation for user behavior, notification performance,
campaign effectiveness, and system-wide metrics.

```mermaid
flowchart TD
    Analytics_Service["Analytics Service"]
    analyticsreportrequest@{ shape: das, label: "analytics.report.request" }
      analyticsreportrequest -->|"Reply"| Analytics_Service
    campaignanalytics@{ shape: das, label: "campaign.analytics" }
      campaignanalytics -->|"Receive"| Analytics_Service
    notificationanalytics@{ shape: das, label: "notification.analytics" }
      notificationanalytics -->|"Receive"| Analytics_Service
    useranalytics@{ shape: das, label: "user.analytics" }
      useranalytics -->|"Receive"| Analytics_Service
    analyticsalert@{ shape: das, label: "analytics.alert" }
      Analytics_Service -->|"Send"| analyticsalert
    analyticsinsights@{ shape: das, label: "analytics.insights" }
      Analytics_Service -->|"Send"| analyticsinsights
    Campaign_Service["Campaign Service"]
          Campaign_Service --> campaignanalytics
    Notification_Service["Notification Service"]
          Notification_Service --> notificationanalytics
    User_Service["User Service"]
          User_Service --> useranalytics

    style Analytics_Service fill:#3498db,stroke:#2980b9,stroke-width:2px,color:#fff
    style Campaign_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Notification_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style User_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

### Campaign Service

A service that manages notification campaigns, user targeting, and campaign execution.
Handles campaign creation, user segmentation, scheduling, and personalized notification delivery.
Uses user data for targeting and personalization of campaign messages.

```mermaid
flowchart TD
    Campaign_Service["Campaign Service"]
    campaigncreate@{ shape: das, label: "campaign.create" }
      campaigncreate -->|"Receive"| Campaign_Service
    campaignexecute@{ shape: das, label: "campaign.execute" }
      campaignexecute -->|"Receive"| Campaign_Service
    campaignanalytics@{ shape: das, label: "campaign.analytics" }
      Campaign_Service -->|"Send"| campaignanalytics
    campaignstatus@{ shape: das, label: "campaign.status" }
      Campaign_Service -->|"Send"| campaignstatus
    notificationuseruser_idpush@{ shape: das, label: "notification.user.{user_id}.push" }
      Campaign_Service -->|"Send"| notificationuseruser_idpush
    userinforequest@{ shape: das, label: "user.info.request" }
      Campaign_Service -->|"Request"| userinforequest
    Analytics_Service["Analytics Service"]
          campaignanalytics --> Analytics_Service
    Notification_Service["Notification Service"]
          notificationuseruser_idpush --> Notification_Service
    User_Service["User Service"]
          userinforequest --> User_Service

    style Campaign_Service fill:#3498db,stroke:#2980b9,stroke-width:2px,color:#fff
    style Analytics_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Notification_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style User_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

### Notification Service

A service that handles user notifications, preferences, and interactions.
Supports real-time notifications, user preferences management.

```mermaid
flowchart TD
    Notification_Service["Notification Service"]
    notificationpreferencesget@{ shape: das, label: "notification.preferences.get" }
      notificationpreferencesget -->|"Reply"| Notification_Service
    notificationpreferencesupdate@{ shape: das, label: "notification.preferences.update" }
      notificationpreferencesupdate -->|"Receive"| Notification_Service
    notificationuseruser_idpush@{ shape: das, label: "notification.user.{user_id}.push" }
      notificationuseruser_idpush -->|"Receive"| Notification_Service
    notificationanalytics@{ shape: das, label: "notification.analytics" }
      Notification_Service -->|"Send"| notificationanalytics
    userinforequest@{ shape: das, label: "user.info.request" }
      Notification_Service -->|"Request"| userinforequest
    Analytics_Service["Analytics Service"]
          notificationanalytics --> Analytics_Service
    Campaign_Service["Campaign Service"]
          Campaign_Service --> notificationuseruser_idpush
    User_Service["User Service"]
          userinforequest --> User_Service
          User_Service --> notificationpreferencesupdate

    style Notification_Service fill:#3498db,stroke:#2980b9,stroke-width:2px,color:#fff
    style Analytics_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Campaign_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style User_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

### User Service

A service that manages user information, profiles, and authentication.
Handles user data requests, profile updates, and user lifecycle events.

```mermaid
flowchart TD
    User_Service["User Service"]
    userinforequest@{ shape: das, label: "user.info.request" }
      userinforequest -->|"Reply"| User_Service
    notificationpreferencesupdate@{ shape: das, label: "notification.preferences.update" }
      User_Service -->|"Send"| notificationpreferencesupdate
    useranalytics@{ shape: das, label: "user.analytics" }
      User_Service -->|"Send"| useranalytics
    userinfoupdate@{ shape: das, label: "user.info.update" }
      User_Service -->|"Send"| userinfoupdate
    Analytics_Service["Analytics Service"]
          useranalytics --> Analytics_Service
    Campaign_Service["Campaign Service"]
          Campaign_Service --> userinforequest
    Notification_Service["Notification Service"]
          notificationpreferencesupdate --> Notification_Service
          Notification_Service --> userinforequest

    style User_Service fill:#3498db,stroke:#2980b9,stroke-width:2px,color:#fff
    style Analytics_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Campaign_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Notification_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

## Channels

### analytics.alert
```mermaid
flowchart LR
    analyticsalert@{ shape: das, label: "analytics.alert" }
    Analytics_Service["Analytics Service"]
    Analytics_Service -->|"Send"| analyticsalert
    style Analytics_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**AnalyticsAlertMessage**
```json
{
  "actions": [
    "string"
  ],
  "affected_services": [
    "string[enum:user_service,notification_service,campaign_service]"
  ],
  "alert_id": "string[uuid]",
  "alert_type": "string[enum:anomaly_detected,threshold_exceeded,trend_change,system_issue]",
  "created_at": "string[date-time]",
  "current_value": "number",
  "description": "string",
  "metadata": {
    "environment": "string[enum:development,staging,production]",
    "platform": "string[enum:ios,android,web]",
    "source": "string[enum:mobile,web,api]",
    "version": "string"
  },
  "metric": "string",
  "severity": "string[enum:low,medium,high,critical]",
  "threshold": "number",
  "time_window": "string",
  "title": "string"
}
```

### analytics.insights
```mermaid
flowchart LR
    analyticsinsights@{ shape: das, label: "analytics.insights" }
    Analytics_Service["Analytics Service"]
    Analytics_Service -->|"Send"| analyticsinsights
    style Analytics_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**AnalyticsInsightMessage**
```json
{
  "category": "string[enum:user_behavior,notification_performance,campaign_effectiveness,system_health]",
  "confidence": "number[float]",
  "created_at": "string[date-time]",
  "data_points": [
    "object"
  ],
  "description": "string",
  "insight_id": "string[uuid]",
  "insight_type": "string[enum:trend,anomaly,recommendation,alert]",
  "metadata": {
    "environment": "string[enum:development,staging,production]",
    "platform": "string[enum:ios,android,web]",
    "source": "string[enum:mobile,web,api]",
    "version": "string"
  },
  "recommendations": [
    "string"
  ],
  "severity": "string[enum:low,medium,high,critical]",
  "title": "string"
}
```

### analytics.report.request
```mermaid
flowchart LR
    analyticsreportrequest@{ shape: das, label: "analytics.report.request" }
    Analytics_Service["Analytics Service"]
    analyticsreportrequest -->|"Reply"| Analytics_Service
    style Analytics_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**request**: AnalyticsReportRequestMessage
```json
{
  "created_at": "string[date-time]",
  "filters": {
    "campaign_ids": [
      "string[uuid]"
    ],
    "event_types": [
      "string"
    ],
    "user_ids": [
      "string[uuid]"
    ],
    "user_segments": [
      "string[enum:all_users,new_users,active_users,inactive_users,premium_users,free_users]"
    ]
  },
  "format": "string[enum:json,csv,pdf]",
  "metrics": [
    "string[enum:event_count,user_count,conversion_rate,engagement_rate,response_time,error_rate]"
  ],
  "report_id": "string[uuid]",
  "report_type": "string[enum:user_activity,notification_performance,campaign_effectiveness,system_health,custom]",
  "time_range": {
    "end": "string[date-time]",
    "granularity": "string[enum:minute,hour,day,week,month]",
    "start": "string[date-time]"
  }
}
```
**reply**: AnalyticsReportReplyMessage
```json
{
  "data": "object",
  "error": {
    "code": "string",
    "message": "string"
  },
  "generated_at": "string[date-time]",
  "insights": [
    {
      "confidence": "number[float]",
      "data_points": [
        "object"
      ],
      "description": "string",
      "impact": "string[enum:low,medium,high]",
      "title": "string",
      "type": "string[enum:trend,anomaly,correlation,recommendation]"
    }
  ],
  "report_id": "string[uuid]",
  "report_type": "string[enum:user_activity,notification_performance,campaign_effectiveness,system_health,custom]",
  "summary": {
    "event_types": "object",
    "top_metrics": {
      "conversion_rate": "number[float]",
      "engagement_rate": "number[float]",
      "error_rate": "number[float]",
      "response_time_avg": "number[float]"
    },
    "total_events": "integer",
    "unique_users": "integer"
  },
  "time_range": {
    "end": "string[date-time]",
    "granularity": "string[enum:minute,hour,day,week,month]",
    "start": "string[date-time]"
  }
}
```

### campaign.analytics
```mermaid
flowchart LR
    campaignanalytics@{ shape: das, label: "campaign.analytics" }
    Campaign_Service["Campaign Service"]
    Campaign_Service -->|"Send"| campaignanalytics
    Analytics_Service["Analytics Service"]
    campaignanalytics -->|"Receive"| Analytics_Service
    style Campaign_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Analytics_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**CampaignAnalyticsEventMessage**
```json
{
  "campaign_id": "string[uuid]",
  "event_id": "string[uuid]",
  "event_type": "string[enum:campaign_created,campaign_executed,notification_sent,notification_opened,notification_clicked,campaign_completed,campaign_failed]",
  "execution_id": "string[uuid]",
  "metadata": {
    "environment": "string[enum:development,staging,production]",
    "platform": "string[enum:ios,android,web]",
    "source": "string[enum:mobile,web,api]",
    "version": "string"
  },
  "notification_id": "string[uuid]",
  "timestamp": "string[date-time]",
  "user_id": "string[uuid]"
}
```

### campaign.create
```mermaid
flowchart LR
    campaigncreate@{ shape: das, label: "campaign.create" }
    Campaign_Service["Campaign Service"]
    campaigncreate -->|"Receive"| Campaign_Service
    style Campaign_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**CampaignCreateMessage**
```json
{
  "campaign_id": "string[uuid]",
  "created_at": "string[date-time]",
  "description": "string",
  "metadata": {
    "environment": "string[enum:development,staging,production]",
    "platform": "string[enum:ios,android,web]",
    "source": "string[enum:mobile,web,api]",
    "version": "string"
  },
  "name": "string",
  "notification_template": {
    "body_template": "string",
    "data": "object",
    "localization": "object",
    "priority": "string[enum:low,normal,high]",
    "title_template": "string"
  },
  "schedule": {
    "recurring": {
      "end_date": "string[date]",
      "frequency": "string[enum:daily,weekly,monthly]",
      "interval": "integer",
      "start_date": "string[date]"
    },
    "scheduled_at": "string[date-time]",
    "timezone": "string",
    "type": "string[enum:immediate,scheduled,recurring]"
  },
  "settings": {
    "a_b_testing": {
      "enabled": "boolean",
      "traffic_split": [
        "number"
      ],
      "variants": [
        {
          "body_template": "string",
          "data": "object",
          "localization": "object",
          "priority": "string[enum:low,normal,high]",
          "title_template": "string"
        }
      ]
    },
    "batch_size": "integer",
    "max_retries": "integer",
    "rate_limit": "integer",
    "respect_quiet_hours": "boolean"
  },
  "target_audience": {
    "estimated_reach": "integer",
    "user_filters": {
      "language": [
        "string"
      ],
      "last_activity": {
        "from": "string[date-time]",
        "to": "string[date-time]"
      },
      "registration_date": {
        "from": "string[date]",
        "to": "string[date]"
      },
      "timezone": [
        "string"
      ]
    },
    "user_segments": [
      "string[enum:all_users,new_users,active_users,inactive_users,premium_users,free_users]"
    ]
  }
}
```

### campaign.execute
```mermaid
flowchart LR
    campaignexecute@{ shape: das, label: "campaign.execute" }
    Campaign_Service["Campaign Service"]
    campaignexecute -->|"Receive"| Campaign_Service
    style Campaign_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**CampaignExecuteMessage**
```json
{
  "batch_size": "integer",
  "campaign_id": "string[uuid]",
  "created_at": "string[date-time]",
  "execution_id": "string[uuid]",
  "execution_type": "string[enum:immediate,scheduled,batch]",
  "metadata": {
    "environment": "string[enum:development,staging,production]",
    "platform": "string[enum:ios,android,web]",
    "source": "string[enum:mobile,web,api]",
    "version": "string"
  },
  "priority": "string[enum:low,normal,high]"
}
```

### campaign.status
```mermaid
flowchart LR
    campaignstatus@{ shape: das, label: "campaign.status" }
    Campaign_Service["Campaign Service"]
    Campaign_Service -->|"Send"| campaignstatus
    style Campaign_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**CampaignStatusUpdateMessage**
```json
{
  "campaign_id": "string[uuid]",
  "error": {
    "code": "string",
    "message": "string"
  },
  "execution_id": "string[uuid]",
  "progress": {
    "failed": "integer",
    "sent": "integer",
    "success_rate": "number[float]",
    "total_targets": "integer"
  },
  "status": "string[enum:pending,running,completed,failed,paused,cancelled]",
  "updated_at": "string[date-time]"
}
```

### notification.analytics
```mermaid
flowchart LR
    notificationanalytics@{ shape: das, label: "notification.analytics" }
    Notification_Service["Notification Service"]
    Notification_Service -->|"Send"| notificationanalytics
    Analytics_Service["Analytics Service"]
    notificationanalytics -->|"Receive"| Analytics_Service
    style Notification_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Analytics_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**NotificationAnalyticsEventMessage**
```json
{
  "event_id": "string[uuid]",
  "event_type": "string[enum:notification_sent,notification_opened,notification_clicked]",
  "metadata": {
    "environment": "string[enum:development,staging,production]",
    "platform": "string[enum:ios,android,web]",
    "source": "string[enum:mobile,web,api]",
    "version": "string"
  },
  "notification_id": "string[uuid]",
  "timestamp": "string[date-time]",
  "user_id": "string[uuid]"
}
```

### notification.preferences.get
```mermaid
flowchart LR
    notificationpreferencesget@{ shape: das, label: "notification.preferences.get" }
    Notification_Service["Notification Service"]
    notificationpreferencesget -->|"Reply"| Notification_Service
    style Notification_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**request**: PreferencesRequestMessage
```json
{
  "user_id": "string[uuid]"
}
```
**reply**: PreferencesReplyMessage
```json
{
  "preferences": {
    "categories": {
      "marketing": "boolean",
      "security": "boolean",
      "updates": "boolean"
    },
    "email_enabled": "boolean",
    "push_enabled": "boolean",
    "quiet_hours": {
      "enabled": "boolean",
      "end": "string[time]",
      "start": "string[time]"
    },
    "sms_enabled": "boolean"
  },
  "updated_at": "string[date-time]"
}
```

### notification.preferences.update
```mermaid
flowchart LR
    notificationpreferencesupdate@{ shape: das, label: "notification.preferences.update" }
    User_Service["User Service"]
    User_Service -->|"Send"| notificationpreferencesupdate
    Notification_Service["Notification Service"]
    notificationpreferencesupdate -->|"Receive"| Notification_Service
    style User_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Notification_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**PreferencesUpdateMessage**
```json
{
  "preferences": {
    "categories": {
      "marketing": "boolean",
      "security": "boolean",
      "updates": "boolean"
    },
    "email_enabled": "boolean",
    "push_enabled": "boolean",
    "quiet_hours": {
      "enabled": "boolean",
      "end": "string[time]",
      "start": "string[time]"
    },
    "sms_enabled": "boolean"
  },
  "updated_at": "string[date-time]",
  "user_id": "string[uuid]"
}
```

### notification.user.{user_id}.push
```mermaid
flowchart LR
    notificationuseruser_idpush@{ shape: das, label: "notification.user.{user_id}.push" }
    Campaign_Service["Campaign Service"]
    Campaign_Service -->|"Send"| notificationuseruser_idpush
    Notification_Service["Notification Service"]
    notificationuseruser_idpush -->|"Receive"| Notification_Service
    style Campaign_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Notification_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**PushNotificationMessage**
```json
{
  "body": "string",
  "created_at": "string[date-time]",
  "data": "object",
  "notification_id": "string[uuid]",
  "priority": "string[enum:low,normal,high]",
  "title": "string",
  "user_id": "string[uuid]"
}
```

### user.analytics
```mermaid
flowchart LR
    useranalytics@{ shape: das, label: "user.analytics" }
    User_Service["User Service"]
    User_Service -->|"Send"| useranalytics
    Analytics_Service["Analytics Service"]
    useranalytics -->|"Receive"| Analytics_Service
    style User_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Analytics_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**UserAnalyticsEventMessage**
```json
{
  "event_id": "string[uuid]",
  "event_type": "string[enum:user_registered,user_logged_in,profile_updated,preferences_changed,account_deleted]",
  "metadata": {
    "environment": "string[enum:development,staging,production]",
    "platform": "string[enum:ios,android,web]",
    "source": "string[enum:mobile,web,api]",
    "version": "string"
  },
  "timestamp": "string[date-time]",
  "user_id": "string[uuid]"
}
```

### user.info.request
```mermaid
flowchart LR
    userinforequest@{ shape: das, label: "user.info.request" }
    Campaign_Service["Campaign Service"]
    Campaign_Service -->|"Request"| userinforequest
    Notification_Service["Notification Service"]
    Notification_Service -->|"Request"| userinforequest
    User_Service["User Service"]
    userinforequest -->|"Reply"| User_Service
    style Campaign_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style Notification_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff
    style User_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**request**: UserInfoRequestMessage
```json
{
  "user_id": "string[uuid]"
}
```
**reply**: UserInfoReplyMessage
```json
{
  "email": "string[email]",
  "error": {
    "code": "string",
    "message": "string"
  },
  "language": "string",
  "name": "string",
  "timezone": "string",
  "user_id": "string[uuid]"
}
```

### user.info.update
```mermaid
flowchart LR
    userinfoupdate@{ shape: das, label: "user.info.update" }
    User_Service["User Service"]
    User_Service -->|"Send"| userinfoupdate
    style User_Service fill:#95A5A6,stroke:#7F8C8D,stroke-width:2px,color:#fff

```

#### Messages
**UserInfoUpdateMessage**
```json
{
  "changes": "object",
  "metadata": {
    "environment": "string[enum:development,staging,production]",
    "platform": "string[enum:ios,android,web]",
    "source": "string[enum:mobile,web,api]",
    "version": "string"
  },
  "updated_at": "string[date-time]",
  "user_id": "string[uuid]"
}
```
