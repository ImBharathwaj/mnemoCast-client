#!/bin/bash

# MnemoCast Screen System Status Check Script

echo "🔍 Checking MnemoCast Screen System Status..."
echo "=============================================="
echo ""

# Check process
echo "📋 Process Status:"
if pgrep -f "bin/screen" > /dev/null; then
    echo "   ✅ RUNNING"
    PROCESS_INFO=$(ps aux | grep "[b]in/screen" | head -1)
    PID=$(echo $PROCESS_INFO | awk '{print $2}')
    CPU=$(echo $PROCESS_INFO | awk '{print $3}')
    MEM=$(echo $PROCESS_INFO | awk '{print $4}')
    RUNTIME=$(ps -o etime= -p $PID 2>/dev/null | tr -d ' ')
    echo "   PID: $PID"
    echo "   CPU: ${CPU}%"
    echo "   Memory: ${MEM}%"
    echo "   Runtime: $RUNTIME"
else
    echo "   ❌ NOT RUNNING"
fi

echo ""

# Check configuration
echo "⚙️  Configuration Status:"
if [ -f ~/.mnemocast/identity.json ]; then
    echo "   ✅ Identity file exists"
    if command -v jq &> /dev/null; then
        SCREEN_ID=$(jq -r '.id' ~/.mnemocast/identity.json 2>/dev/null)
        SCREEN_NAME=$(jq -r '.name' ~/.mnemocast/identity.json 2>/dev/null)
        echo "   Screen ID: $SCREEN_ID"
        echo "   Screen Name: $SCREEN_NAME"
    else
        SCREEN_ID=$(grep -o '"id":"[^"]*' ~/.mnemocast/identity.json | cut -d'"' -f4)
        echo "   Screen ID: $SCREEN_ID"
    fi
else
    echo "   ⚠️  Identity file NOT FOUND"
fi

if [ -f ~/.mnemocast/config.json ]; then
    echo "   ✅ Config file exists"
    if command -v jq &> /dev/null; then
        AD_SERVER=$(jq -r '.adServerUrl' ~/.mnemocast/config.json 2>/dev/null)
        HEARTBEAT_INT=$(jq -r '.heartbeatInterval' ~/.mnemocast/config.json 2>/dev/null)
        echo "   Ad Server: $AD_SERVER"
        echo "   Heartbeat Interval: ${HEARTBEAT_INT}s"
    fi
else
    echo "   ⚠️  Config file NOT FOUND"
fi

echo ""

# Check credentials
echo "🔐 Credentials Status:"
if [ -f ~/.mnemocast/credentials.json.enc ]; then
    echo "   ✅ Credentials file exists (encrypted)"
    if [ -f ~/.mnemocast/.encryption_key ]; then
        echo "   ✅ Encryption key exists"
    else
        echo "   ⚠️  Encryption key NOT FOUND"
    fi
else
    echo "   ⚠️  Credentials NOT CONFIGURED"
fi

echo ""

# Check backend connectivity
echo "🌐 Backend Connectivity:"
BACKEND_URL="http://10.42.0.1:8080"
if command -v curl &> /dev/null; then
    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" --connect-timeout 3 "$BACKEND_URL/api/v1/health" 2>/dev/null)
    if [ "$HTTP_CODE" = "200" ]; then
        echo "   ✅ Backend is ACCESSIBLE"
        echo "   URL: $BACKEND_URL"
    else
        echo "   ❌ Backend is NOT ACCESSIBLE (HTTP $HTTP_CODE)"
        echo "   URL: $BACKEND_URL"
    fi
else
    echo "   ⚠️  curl not available - cannot check backend"
fi

echo ""

# Check if screen is registered (if we have screen ID)
if [ -n "$SCREEN_ID" ] && command -v curl &> /dev/null; then
    echo "📡 Registration Status:"
    REG_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BACKEND_URL/api/v1/screens/$SCREEN_ID" 2>/dev/null)
    if [ "$REG_RESPONSE" = "200" ]; then
        echo "   ✅ Screen is REGISTERED in backend"
    elif [ "$REG_RESPONSE" = "404" ]; then
        echo "   ⚠️  Screen is NOT REGISTERED in backend"
    else
        echo "   ❓ Registration status unknown (HTTP $REG_RESPONSE)"
    fi
fi

echo ""
echo "=============================================="
echo "Status check complete!"
echo ""

