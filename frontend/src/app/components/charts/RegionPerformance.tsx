import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/app/components/ui/card"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

interface RegionPerformanceProps {
  data: { region: string; kills: number; isk: number }[]
}

export default function RegionPerformance({ data }: RegionPerformanceProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Region Performance</CardTitle>
        <CardDescription>Comparison of kills and ISK destroyed by region</CardDescription>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={400}>
          <BarChart data={data}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="region" />
            <YAxis yAxisId="left" orientation="left" stroke="#8884d8" />
            <YAxis yAxisId="right" orientation="right" stroke="#82ca9d" />
            <Tooltip />
            <Bar yAxisId="left" dataKey="kills" fill="#8884d8" name="Kills" />
            <Bar yAxisId="right" dataKey="isk" fill="#82ca9d" name="ISK Destroyed" />
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}