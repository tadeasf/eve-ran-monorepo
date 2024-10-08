import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/app/components/ui/card"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts'

interface CharacterPerformanceProps {
  data: { name: string; kills: number; isk: number }[]
}

export default function CharacterPerformance({ data }: CharacterPerformanceProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Character Performance</CardTitle>
        <CardDescription>Top performers by kills and ISK destroyed</CardDescription>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={400}>
          <BarChart data={data} layout="vertical">
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis type="number" />
            <YAxis dataKey="name" type="category" />
            <Tooltip />
            <Bar dataKey="kills" fill="#8884d8" name="Kills" />
            <Bar dataKey="isk" fill="#82ca9d" name="ISK Destroyed" />
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}