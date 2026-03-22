import React from 'react';
import { Box, Typography } from '@mui/material';

interface StatBoxProps {
  label: string;
  value: number | null;
  unit?: string;
}

function StatBox({ label, value, unit }: StatBoxProps) {
  return (
    <Box sx={{ textAlign: 'center', px: 1 }}>
      <Typography variant="h5" fontWeight="bold" color="primary">
        {value !== null ? value.toFixed(1) : '—'}
        {value !== null && unit ? (
          <Typography component="span" variant="body2" color="text.secondary" sx={{ ml: 0.5 }}>
            {unit}
          </Typography>
        ) : null}
      </Typography>
      <Typography variant="caption" color="text.secondary" sx={{ textTransform: 'uppercase', letterSpacing: 0.5 }}>
        {label}
      </Typography>
    </Box>
  );
}

export default StatBox;
