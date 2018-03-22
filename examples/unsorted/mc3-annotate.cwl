cwlVersion: v1.0
class: Workflow

requirements:
  - class: SubworkflowFeatureRequirement
  - class: ScatterFeatureRequirement
  - class: StepInputExpressionRequirement
  - class: InlineJavascriptRequirement
  - class: MultipleInputFeatureRequirement

inputs:
  tumor:
    type: File
    secondaryFiles:
      - .bai
  normal:
    type: File
    secondaryFiles:
      - .bai
  reference:
    type: File
    secondaryFiles:
      - .fai

  reference_dict:
    type: File

  dbsnp:
    type: File
    secondaryFiles:
      - .tbi
  cosmic:
    type: File
    secondaryFiles:
      - .tbi
  centromere:
    type: File

  tumor_id:
    type: string
  normal_id:
    type: string

  # TODO workaround for unknown tn ids
  somaticsniper_normal_id:
    type: string
  somaticsniper_tumor_id:
    type: string
  varscan_snp_normal_id:
    type: string
  varscan_snp_tumor_id:
    type: string
  varscan_indel_normal_id:
    type: string
  varscan_indel_tumor_id:
    type: string

  #genome_reference:
  #  type: File

  #cosmic_reference:
  #  type: File

  #valstatus_database:
  #  type: File

outputs:
  a:
    type: File
    outputSource: normalize_muse/output_vcf

  b:
    type: File
    outputSource: normalize_radia/output_vcf

  c:
    type: File
    outputSource: normalize_somaticsniper/output_vcf

  d:
    type: File
    outputSource: normalize_varscan_snp/output_vcf

  e:
    type: File
    outputSource: normalize_varscan_indel/output_vcf

  f:
    type: File
    outputSource: merge/output_vcf

  #output_maf:
    #type: File
    #outputSource: convert_vcf_to_maf/output_maf

steps:
  mc3:
    in: 
      tumor: tumor
      normal: normal
      reference: reference
      dbsnp: dbsnp
      cosmic: cosmic
      centromere: centromere
    run: mc3/workflows/mc3_full.cwl
    out:
      - pindel-out
      - somaticsniper-out
      - varscan-snp-out
      - varscan-indel-out
      - muse-out
      - mutect-out
      - radia-out

  filter_muse:
    in:
      vcf: mc3/muse-out
      output_name:
        valueFrom: muse.filtered.vcf
    out: [output_vcf]
    run: tools/filter_muse/tool.cwl.yaml

  sort_muse:
    in:
      seqdict: reference_dict
      vcf: filter_muse/output_vcf
      output_name:
        valueFrom: muse.sorted.vcf
    out: [output_vcf]
    run: tools/sort_vcf/tool.cwl.yaml

  normalize_muse:
    in:
      vcf: sort_muse/output_vcf
      output_name:
        valueFrom: muse.normalized.vcf
      vcf_normal_id: normal_id
      vcf_tumor_id: tumor_id
    out: [output_vcf]
    run: tools/vcf2maf/vcf2vcf.cwl.yaml

  filter_radia:
    in:
      vcf: mc3/radia-out
      output_name:
        valueFrom: radia.filtered.vcf
    out: [output_vcf]
    run: tools/filter_radia/tool.cwl.yaml

  sort_radia:
    in:
      seqdict: reference_dict
      vcf: filter_radia/output_vcf
      output_name:
        valueFrom: radia.sorted.vcf
    out: [output_vcf]
    run: tools/sort_vcf/tool.cwl.yaml

  normalize_radia:
    in:
      vcf: sort_radia/output_vcf
      output_name:
        valueFrom: radia.normalized.vcf
      vcf_normal_id: normal_id
      vcf_tumor_id: tumor_id
    out: [output_vcf]
    run: tools/vcf2maf/vcf2vcf.cwl.yaml

  filter_somaticsniper:
    in:
      vcf: mc3/somaticsniper-out
      normal_id: normal_id
      tumor_id: tumor_id
      output_name:
        valueFrom: somaticsniper.filtered.vcf
    out: [output_vcf]
    run: tools/filter_somaticsniper/tool.cwl.yaml

  sort_somaticsniper:
    in:
      seqdict: reference_dict
      vcf: filter_somaticsniper/output_vcf
      output_name:
        valueFrom: somaticsniper.sorted.vcf
    out: [output_vcf]
    run: tools/sort_vcf/tool.cwl.yaml

  normalize_somaticsniper:
    in:
      vcf: sort_somaticsniper/output_vcf
      output_name:
        valueFrom: somaticsniper.normalized.vcf
      vcf_normal_id: somaticsniper_normal_id
      vcf_tumor_id: somaticsniper_tumor_id
    out: [output_vcf]
    run: tools/vcf2maf/vcf2vcf.cwl.yaml

#  preprocess_muse:
#    in:
#      seqdict: reference_dict
#      vcf: mc3/muse-out
#    out: [output_vcf]
#    run: tools/preprocess_muse/tool.cwl.yaml

#  preprocess_radia:
#    in:
#      seqdict: reference_dict
#      vcf: mc3/radia-out
#    out: [output_vcf]
#    run: tools/preprocess_radia/tool.cwl.yaml

#  preprocess_somaticsniper:
#    in:
#      vcf: mc3/somaticsniper-out
#      # TODO get_tn_ids is broken
#      #normal_id: normal_id
#      #tumor_id: tumor_id
#      seqdict: reference_dict
#    out: [output_vcf]
#    run: tools/preprocess_somaticsniper/tool.cwl.yaml

  filter_varscan_snp:
    in:
      vcf: mc3/varscan-snp-out
      # TODO get_tn_ids is broken
      #normal_id: varscan_snp_tn_ids/normal
      #tumor_id: varscan_snp_tn_ids/tumor
      normal_id: varscan_snp_normal_id
      tumor_id: varscan_snp_tumor_id
      seqdict: reference_dict
    out: [output_vcf]
    run: tools/vcf2vcf_filter/tool.cwl.yaml

  sort_varscan_snp:
    in:
      seqdict: reference_dict
      vcf: filter_varscan_snp/output_vcf
      output_name:
        valueFrom: varscan.snp.sorted.vcf
    out: [output_vcf]
    run: tools/sort_vcf/tool.cwl.yaml

#  varscan_snp_tn_ids:
#    in:
#      vcf: mc3/varscan-snp-out
#    out:
#      - tumor
#      - normal
#    run: tools/get_tn_ids/tool.cwl.yaml
#
#  varscan_indel_tn_ids:
#    in:
#      vcf: mc3/varscan-indel-out
#    out:
#      - tumor
#      - normal
#    run: tools/get_tn_ids/tool.cwl.yaml
#
#  somaticsniper_tn_ids:
#    in:
#      vcf: mc3/somaticsniper-out
#    out:
#      - tumor
#      - normal
#    run: tools/get_tn_ids/tool.cwl.yaml

  normalize_varscan_snp:
    in:
      vcf: sort_varscan_snp/output_vcf
      output_name:
        valueFrom: varscan.snp.normalized.vcf
      vcf_tumor_id: tumor_id
      vcf_normal_id: normal_id
    out: [output_vcf]
    run: tools/vcf2maf/vcf2vcf.cwl.yaml

  filter_varscan_indel:
    in:
      vcf: mc3/varscan-indel-out
      # TODO get_tn_ids is broken
      #normal_id: varscan_indel_tn_ids/normal
      #tumor_id: varscan_indel_tn_ids/tumor
      normal_id: varscan_indel_normal_id
      tumor_id: varscan_indel_tumor_id
      seqdict: reference_dict
    out: [output_vcf]
    run: tools/vcf2vcf_filter/tool.cwl.yaml

  sort_varscan_indel:
    in:
      seqdict: reference_dict
      vcf: filter_varscan_indel/output_vcf
      output_name:
        valueFrom: varscan.indel.sorted.vcf
    out: [output_vcf]
    run: tools/sort_vcf/tool.cwl.yaml

  normalize_varscan_indel:
    in:
      vcf: sort_varscan_indel/output_vcf
      output_name:
        valueFrom: varscan.indel.normalized.vcf
      vcf_tumor_id: tumor_id
      vcf_normal_id: normal_id
    out: [output_vcf]
    run: tools/vcf2maf/vcf2vcf.cwl.yaml

  merge:
    in:
      keys:
        default:
          - muse
          - radia
          - somaticsniper
          - varscan-snp
          - varscan-indel
      vcfs:
        - normalize_muse/output_vcf
        - normalize_radia/output_vcf
        - normalize_somaticsniper/output_vcf
        - normalize_varscan_snp/output_vcf
        - normalize_varscan_indel/output_vcf
    out: [output_vcf]
    run: tools/merge_vcfs/tool.cwl.yaml

#  variant_effect_predictor:
#    label: Variant effect predictor (VEP)
#    doc:
#    in:
#      vcf: merge_vcfs/output_vcf
#      keys: keys
#    out: [output_vcf]
#    run: ../variant_effect_predictor/tool.cwl.yaml
#
#  annotate_vcf_cosmic:
#    label: Annotate VCF Cosmic
#    doc:
#    in:
#      genome_reference: genome_reference
#      cosmic_reference: cosmic_reference
#      valstatus_database: valstatus_database
#      vcf: variant_effect_predictor/output_vcf
#    out: [output_vcf]
#    run: ../annotate_vcf_cosmic/tool.cwl.yaml
#
#  convert_vcf_to_maf:
#    label: VCF to MAF
#    doc:
#    in:
#      vcf: annotate_vcf_cosmic/output_vcf
#      tumor_id: tumor_id
#      normal_id: normal_id
#    out: [output_maf]
#    run: ../convert_vcf_to_maf/tool.cwl.yaml
