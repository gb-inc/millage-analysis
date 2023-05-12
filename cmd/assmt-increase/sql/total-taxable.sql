SELECT SUM(AN.ImprovementAssessment + AN.LandAssessment) TotalTxblAmt FROM
Property P
INNER JOIN AssmtNew AN
ON P.PropertyID = AN.PropertyID
INNER JOIN TownshipBoroughCodes tbc
on P.TownShipBorough = tbc.TownshipID
INNER JOIN LandUseCode luc
ON P.LandUse = luc.Code
WHERE An.CurrentAssessment <> 0
AND P.ParcelEliminated IS NULL
AND P.LandUse NOT IN ('166', 'W10')
AND tbc.TownshipOrBorough IN ('T','B')
AND NOT luc.ParcelStatus = 'N'