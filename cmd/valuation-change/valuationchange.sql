SELECT
 P.TownShipBorough
,SUM(ANOld.LandAssessment) AS OldLandAssmt
,SUM(ANOld.ImprovementAssessment) AS OldImprAssmt
,SUM(AN.LandAssessment) AS NewLandAssmt
,SUM(AN.ImprovementAssessment) AS NewImprAssmt
,SUM((AN.LandAssessment - COALESCE(ANOLD.LandAssessment, 0.00))) LandDiff
,SUM((AN.ImprovementAssessment - COALESCE(ANOLD.ImprovementAssessment, 0.00))) ImprDiff
FROM AssmtNew AN
INNER JOIN Property P
    ON P.PropertyID = AN.PropertyID
INNER JOIN TownshipBoroughCodes TBC
    ON P.TownShipBorough = tbc.TownshipID
INNER JOIN LandUseCode LUC
    ON P.LandUse = LUC.Code
LEFT JOIN (SELECT * FROM AssmtNew WHERE EffectiveDate <= '2023-01-01 00:00:00.000' AND CurrentAssessment = -1) ANOLD
    ON P.PropertyID = ANOld.PropertyID
WHERE AN.CurrentAssessment = -1
AND P.ParcelEliminated IS NULL
AND P.LandUse NOT IN ('166','W10')
AND TBC.TownshipOrBorough IN ('T','B')
AND NOT luc.ParcelStatus = 'N'
/* Wayne Highlands School District Specific */
AND P.TownShipBorough IN
('030',
 '070',
 '150',
 '200',
 '230',
 '020',
 '090',
 '110',
 '130',
 '010',
 '050',
 '170',
 '210',
 '270') -- 273 in Wallenpaupack
/* Wayne Highlands School District Specific */
GROUP BY TownShipBorough
ORDER BY TownShipBorough
