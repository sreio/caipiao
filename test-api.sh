#!/bin/bash

echo "========================================="
echo "  测试官方彩票API连接"
echo "========================================="
echo ""

# 测试双色球API
echo "📍 测试双色球API..."
echo "URL: https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=ssq&issueCount=1"
echo ""

ssq_response=$(curl -s "https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice?name=ssq&issueCount=1")

if echo "$ssq_response" | grep -q "state"; then
    echo "✅ 双色球API连接成功！"
    echo "$ssq_response" | grep -o '"code":"[^"]*"' | head -1 || echo ""
else
    echo "❌ 双色球API连接失败"
fi

echo ""
echo "---"
echo ""

# 测试大乐透API
echo "📍 测试大乐透API..."
echo "URL: https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=85&provinceId=0&pageSize=1&isVerify=1&pageNo=1"
echo ""

dlt_response=$(curl -s "https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry?gameNo=85&provinceId=0&pageSize=1&isVerify=1&pageNo=1")

if echo "$dlt_response" | grep -q "success"; then
    echo "✅ 大乐透API连接成功！"
    echo "$dlt_response" | grep -o '"lotteryDrawNum":"[^"]*"' | head -1 || echo ""
else
    echo "❌ 大乐透API连接失败"
fi

echo ""
echo "========================================="
echo "  测试完成"
echo "========================================="
echo ""

echo "💡 如果API测试失败，可能的原因："
echo "   1. 网络连接问题"
echo "   2. 防火墙或代理阻止"
echo "   3. 官方API维护中"
echo ""
echo "💡 如果测试成功，可以直接启动项目使用！"
echo ""

