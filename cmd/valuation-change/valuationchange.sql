SELECT
 P.TownShipBorough
--,AN.PropertyID
--,AN.AssessmentID
--,AN.PriorAssmtID
--,AN.EffectiveDate
,SUM(ANOld.LandAssessment) AS OldLandAssmt
,SUM(ANOld.ImprovementAssessment) AS OldImprAssmt
,SUM(AN.LandAssessment) AS NewLandAssmt
,SUM(AN.ImprovementAssessment) AS NewImprAssmt
,SUM((AN.LandAssessment - ANOld.LandAssessment)) LandDiff
,SUM((AN.ImprovementAssessment - ANOLD.ImprovementAssessment)) ImprDiff
FROM AssmtNew AN
INNER JOIN Property P
    ON P.PropertyID = AN.PropertyID
LEFT JOIN (SELECT * FROM AssmtNew WHERE EffectiveDate = '2023-01-01 00:00:00.000') ANOld
    ON ANOld.PropertyID = AN.PropertyID
WHERE AN.CurrentAssessment = -1
AND AN.EffectiveDate > '2023-01-01 00:00:00.000'
AND (AN.LandAssessment <> ANOld.LandAssessment OR AN.ImprovementAssessment <> ANOld.ImprovementAssessment) --Detect at least one difference in assmt
AND P.ParcelEliminated IS NULL
GROUP BY TownShipBorough
ORDER BY TownShipBorough